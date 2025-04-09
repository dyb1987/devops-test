package service

import (
	"context"
	"exchange-jenkins-cli/config"
	"exchange-jenkins-cli/internal/common"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/rs/zerolog"
	"os"
	"sort"
	"strings"
	"time"
)

type Apps struct {
	jenkins *gojenkins.Jenkins
	Logger  zerolog.Logger
	Items   []*App
}

func New() *Apps {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	conf := config.C
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jenkins, err := gojenkins.CreateJenkins(nil, conf.URL, conf.Username, conf.Password).Init(ctx)
	if err != nil {
		logger.Error().Err(err).Send()
		os.Exit(1)
	}
	return &Apps{
		jenkins: jenkins,
		Logger:  logger,
	}
}

// 创建job
func (apps *Apps) Create(templateFilePath string) error {
	for _, app := range apps.Items {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		configXml, err := app.ParseTemplate(templateFilePath)
		if err != nil {
			return err
		}
		_, err = apps.jenkins.CreateJob(ctx, configXml, fmt.Sprintf("%s-%s", app.ProjectName, app.JobName))
		if err != nil {
			apps.Logger.Error().Err(err).Send()
			continue
		}
		apps.Logger.Info().Msg(fmt.Sprintf("%s-%s create success", app.ProjectName, app.JobName))
	}
	return nil
}

// 启动项目，按照顺序自动化构建
func (apps *Apps) Run() {
	// 1. 原地排序
	sort.Sort(sort.Reverse(AppSlice(apps.Items)))

	// 2. 构建
	for _, app := range apps.Items {
		takes, err := apps.build(fmt.Sprintf("%s-%s", app.ProjectName, app.JobName), map[string]string{"BRANCH": app.BranchName})
		if err != nil {
			apps.Logger.Error().Msg(fmt.Sprintf("%s-%s build failed", app.ProjectName, app.JobName))
			continue
		}
		apps.Logger.Info().Msg(fmt.Sprintf("%s-%s build success, takes %.2f seconds", app.ProjectName, app.JobName, takes))
	}
}

// 构建job
func (apps *Apps) build(jobName string, params map[string]string) (float64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 构建job
	queueId, err := apps.jenkins.BuildJob(ctx, jobName, params)
	if err != nil {
		apps.Logger.Error().Err(err).Send()
	}

	// 获取构建信息
	build, err := apps.jenkins.GetBuildFromQueueID(ctx, queueId)
	if err != nil {
		apps.Logger.Error().Err(err).Send()
	}

	// 轮询job是否构建成功
	for {
		time.Sleep(3 * time.Second)
		if build.IsGood(ctx) {
			// 毫秒
			return build.GetDuration() * 0.001, nil
		}
		time.Sleep(3 * time.Second)
		if !build.IsRunning(ctx) && build.Raw.Result != "SUCCESS" {
			return 0, common.ErrBuildFailed
		}
	}

}

func (apps *Apps) Update(templateFilePath string, names ...string) error {
	if len(names) == 0 {
		return apps.updateAll(templateFilePath)
	}
	var myApps []*App
	for _, name := range names {
		for _, app := range apps.Items {
			if name == fmt.Sprintf("%s-%s", app.ProjectName, app.JobName) {
				myApps = append(myApps, app)
			}
		}
	}
	if len(myApps) == 0 {
		apps.Logger.Error().Msg("the jenkins-job is not exists.")
		return nil
	}

	for _, app := range myApps {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		configXml, err := app.ParseTemplate(templateFilePath)
		if err != nil {
			return err
		}

		err = apps.update(ctx, fmt.Sprintf("%s-%s", app.ProjectName, app.JobName), configXml)
		if err != nil {
			continue
		}
	}
	return nil
}

func (apps *Apps) updateAll(templateFilePath string) error {
	for _, app := range apps.Items {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		configXml, err := app.ParseTemplate(templateFilePath)
		if err != nil {
			return err
		}

		err = apps.update(ctx, fmt.Sprintf("%s-%s", app.ProjectName, app.JobName), configXml)
		if err != nil {
			continue
		}
	}
	return nil
}

func (apps *Apps) update(ctx context.Context, jobName string, configXml string) error {
	job, err := apps.jenkins.GetJob(ctx, jobName)
	if err != nil {
		apps.Logger.Error().Err(err).Msg("get job failed")
		return err
	}
	err = job.UpdateConfig(ctx, configXml)
	if err != nil {
		apps.Logger.Error().Err(err).Msg("update job config failed")
		return err
	}
	apps.Logger.Info().Msg(fmt.Sprintf("%s update success", jobName))
	return nil
}

func (apps *Apps) Delete(view string, names ...string) error {
	jobs := []string{}
	if view == "" {
		jobs = names
	} else {
		var err error
		jobs, err = apps.getJobsFromView(view)
		if err != nil {
			return err
		}
	}
	for _, job := range jobs {
		isDelete, err := apps.delete(job)
		if err != nil {
			apps.Logger.Error().Err(err).Msg(fmt.Sprintf("delete job %s failed", job))
			continue
		}
		if isDelete {
			apps.Logger.Info().Msg(fmt.Sprintf("delete job %s success", job))
		}
	}
	return nil
}

func (apps *Apps) delete(jobName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return apps.jenkins.DeleteJob(ctx, jobName)
}

func (apps *Apps) getJobsFromView(name string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	view, err := apps.jenkins.GetView(ctx, name)
	if err != nil {
		return nil, err
	}

	jobs := []string{}
	for _, job := range view.GetJobs() {
		jobs = append(jobs, job.Name)
	}
	return jobs, nil
}

func (apps *Apps) getConfigXmlByJob(jobName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	job, err := apps.jenkins.GetJob(ctx, jobName)
	if err != nil {
		return "", err
	}
	data, err := job.GetConfig(ctx)
	if err != nil {
		return "", err
	}
	return data, nil
}

type App struct {
	JobName           string `json:"JobName,omitempty"` // jenkins-job
	BranchName        string `json:"BranchName,omitempty"`
	CodeRepository    string `json:"CodeRepository,omitempty"`
	ConfigRepository  string `json:"ConfigRepository,omitempty"`
	ConfigBranch      string `json:"ConfigBranch,omitempty"`
	AppName           string `json:"AppName,omitempty"`
	ServiceName       string `json:"ServiceName,omitempty"`
	ServicePort       string `json:"ServicePort,omitempty"`
	LimitMemory       string `json:"LimitMemory,omitempty"`
	LimitCpu          string `json:"LimitCpu,omitempty"`
	PriorityClassName string `json:"PriorityClassName,omitempty"`
	RequestsMemory    string `json:"RequestsMemory,omitempty"`
	RequestsCpu       string `json:"RequestsCpu,omitempty"`
	Target            string `json:"Target,omitempty"`
	JavaOpts          string `json:"JavaOpts,omitempty"`
	Replicas          string `json:"Replicas,omitempty"`
	Order             int    `json:"Order,omitempty"` // app的启动顺序，数字越大越先启动
	DomainNames       string `json:"DomainNames,omitempty"`
	APIDomain         string `json:"APIDomain,omitempty"`
	// 如下是公共属性
	ImageHarbor string `json:"ImageHarbor,omitempty"`
	ProjectName string `json:"ProjectName,omitempty"`
	TimeZone    string `json:"TimeZone,omitempty"`
	KubeConfig  string `json:"KubeConfig,omitempty"`
}

type AppSlice []*App

func (a AppSlice) Len() int           { return len(a) }
func (a AppSlice) Less(i, j int) bool { return a[i].Order < a[j].Order }
func (a AppSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// 解析模板文件
func (app *App) ParseTemplate(templatePath string) (string, error) {
	data, err := common.ReadXml(templatePath)
	if err != nil {
		return "", err
	}

	replacements := []string{
		"{{ BRANCH_NAME }}", app.BranchName,
		"{{ CODE_REPOSITORY }}", app.CodeRepository,
		"{{ CONFIG_REPOSITORY }}", app.ConfigRepository,
		"{{ CONFIG_BRANCH }}", app.ConfigBranch,
		"{{ PROJECT_NAME }}", app.ProjectName,
		"{{ IMAGE_HARBOR }}", app.ImageHarbor,
		"{{ APP_NAME }}", app.AppName,
		"{{ SERVICE_NAME }}", app.ServiceName,
		"{{ SERVICE_PORT }}", app.ServicePort,
		"{{ LIMIT_MEMORY }}", app.LimitMemory,
		"{{ LIMIT_CPU }}", app.LimitCpu,
		"{{ REQUESTS_MEMORY }}", app.RequestsMemory,
		"{{ REQUESTS_CPU }}", app.RequestsCpu,
		"{{ PRIORITY_CLASS_NAME }}", app.PriorityClassName,
		"{{ REPLICAS }}", app.Replicas,
		"{{ TARGET }}", app.Target,
		"{{ KUBE_CONFIG }}", app.KubeConfig,
		"{{ TIMEZONE }}", app.TimeZone,
		"{{ JAVA_OPTS }}", app.JavaOpts,
		"{{ DOMAIN_NAMES }}", app.DomainNames,
		"{{ API_DOMAIN }}", app.APIDomain,
	}
	replacer := strings.NewReplacer(replacements...)
	configXml := replacer.Replace(data)
	return configXml, nil
}
