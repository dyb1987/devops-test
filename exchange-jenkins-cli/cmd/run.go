package cmd

import (
	"encoding/json"
	"exchange-jenkins-cli/internal/service"
	"github.com/spf13/cobra"
	"os"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "build jenkins job",
	Long:  "运行项目，按照app的顺序构建job",
	Run: func(cmd *cobra.Command, args []string) {
		apps := service.New()
		// 从json中读取数据，填充到apps.Items字段中
		jsonData, err := os.ReadFile("template/" + InputFile)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(jsonData, &apps.Items)
		if err != nil {
			panic(err)
		}
		// 批量构建job
		apps.Run()
	},
}
