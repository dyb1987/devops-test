package cmd

import (
	"exchange-jenkins-cli/internal/service"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete jenkins job",
	Long:  "批量删除jenkins job",
	Run: func(cmd *cobra.Command, args []string) {
		apps := service.New()
		// 删除job
		err := apps.Delete(View, Names...)
		if err != nil {
			panic(err)
		}
	},
}
