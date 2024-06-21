/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/karlhuang95/praetor/api"
	"github.com/spf13/cobra"
)

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "终端启动",
	Long:  `通过终端方式启动,需要手动指定参数`,
	Run: func(cmd *cobra.Command, args []string) {
		httpAddr, _ := cmd.Flags().GetString("http")
		raftAddr, _ := cmd.Flags().GetString("raft")
		myid, _ := cmd.Flags().GetString("myid")
		cluster, _ := cmd.Flags().GetString("cluster")
		api.Start(httpAddr, raftAddr, myid, cluster)
	},
}

func init() {
	rootCmd.AddCommand(consoleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consoleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consoleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	consoleCmd.PersistentFlags().String("http", "127.0.0.1:7001", "http list addr")
	consoleCmd.PersistentFlags().String("raft", "127.0.0.1:7000", "raft list addr")
	consoleCmd.PersistentFlags().String("myid", "1", "raft idr")
	consoleCmd.PersistentFlags().String("cluster", "1/127.0.0.1:7000,2/127.0.0.1:8000,3/127.0.0.1:9000", "cluster info")
}
