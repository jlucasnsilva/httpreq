package cmd

import (
	"fmt"

	"github.com/jlucasnsilva/httpreq"
	"github.com/spf13/cobra"
)

var echoCmd = &cobra.Command{
	Use:   "echo",
	Short: "Starts a echo HTTP server for testing your request groups",
	Long: `Starts a echo HTTP server for testing your request groups. This
server only prints the request and responds with the request body.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := httpreq.EchoServer(":8080"); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(echoCmd)
}
