package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	InputFile    string
	TemplateFile string
	Names        []string
	View         string
)

var rootCmd = &cobra.Command{
	Use:              "exchange-k8s-pipeline",
	Short:            "交易所Jenkins流水线工具",
	Long:             `交易所Jenkins流水线工具，实现对jenkins-job的增删改查功能`,
	TraverseChildren: true,
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(deleteCmd)
	createCmd.PersistentFlags().StringVarP(&InputFile, "input", "i", "", "参考文件")
	createCmd.PersistentFlags().StringVarP(&TemplateFile, "template", "t", "", "模板文件")
	createCmd.MarkFlagsRequiredTogether("input", "template")
	updateCmd.PersistentFlags().StringVarP(&InputFile, "input", "i", "", "参考文件")
	updateCmd.PersistentFlags().StringVarP(&TemplateFile, "template", "t", "", "模板文件")
	updateCmd.PersistentFlags().StringArrayVarP(&Names, "names", "n", []string{}, "job-name")
	updateCmd.MarkFlagsRequiredTogether("input", "template")
	runCmd.PersistentFlags().StringVarP(&InputFile, "input", "i", "", "参考文件")
	_ = runCmd.MarkPersistentFlagRequired("input")
	deleteCmd.PersistentFlags().StringArrayVarP(&Names, "names", "n", []string{}, "job-name")
	deleteCmd.PersistentFlags().StringVarP(&View, "view", "", "", "view-name")
	deleteCmd.MarkFlagsOneRequired("names", "view")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
