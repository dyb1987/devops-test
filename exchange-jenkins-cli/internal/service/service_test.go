package service_test

import (
	"context"
	"encoding/json"
	"exchange-jenkins-cli/config"
	"fmt"
	"github.com/bndr/gojenkins"
	"log"
	"testing"
	"time"
)

var jenkins *gojenkins.Jenkins
var ctx = context.Background()

func init() {
	conf := config.C
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	jenkins, err = gojenkins.CreateJenkins(nil, conf.URL, conf.Username, conf.Password).Init(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetJobConfigXml(t *testing.T) {
	job, err := jenkins.GetJob(ctx, "nmec-prod-web-ui")
	if err != nil {
		log.Fatal(err)
	}
	data, err := job.GetConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}

func TestGetJob(t *testing.T) {
	job, err := jenkins.GetJob(ctx, "testABC-test-contract-option")
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(job.Raw)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}

func TestGetView(t *testing.T) {
	view, err := jenkins.GetView(ctx, "TEST-TEST")
	if err != nil {
		log.Fatal(err)
	}
	for _, job := range view.GetJobs() {
		fmt.Printf("%#v\n", job)
	}
}
