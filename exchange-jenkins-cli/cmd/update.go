package cmd

import (
	"encoding/json"
	"exchange-jenkins-cli/internal/service"
	"github.com/spf13/cobra"
	"os"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update jenkins job",
	Long:  "根据输入源和模板，批量更新jenkins job",
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
		// 更新job
		err = apps.Update("template/"+TemplateFile, Names...)
		if err != nil {
			panic(err)
		}
	},
}
