/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	s "github.com/Yelsnik/broadcast-server/server"

	"github.com/spf13/cobra"
)

var (
	Port  string
	Host  string
	Start string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var port string
		var server *s.Server
		var err error

		if Port != "" {
			port = Port
			server, err = s.NewServer(port)
			if err != nil {
				log.Fatal("cannot create server", err)
				os.Exit(1)
			}

			server.StartServer()

			// Wait for a SIGINT or SIGTERM signal to gracefully shut down the server
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan

			fmt.Println("Shutting down server...")
			server.Stop()
			fmt.Println("Server stopped.")
		}

		fmt.Println("enter a connected port number")

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	startCmd.PersistentFlags().StringVarP(&Port, "port", "p", "", "port to be used to start the server")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
