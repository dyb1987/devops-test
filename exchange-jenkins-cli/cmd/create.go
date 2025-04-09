package cmd

import (
	"encoding/json"
	"exchange-jenkins-cli/internal/service"
	"github.com/spf13/cobra"
	"os"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create jenkins job",
	Long:  "根据输入的json文件和config.xml模板文件，创建jenkins job",
	Run: func(cmd *cobra.Command, args []string) {
		apps := service.New()
		// 从json中读取数据，反序列化为app对象
		jsonData, err := os.ReadFile("template/" + InputFile)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(jsonData, &apps.Items)
		if err != nil {
			panic(err)
		}
		// 批量构建job
		err = apps.Create("template/" + TemplateFile)
		if err != nil {
			panic(err)
		}
	},
}
