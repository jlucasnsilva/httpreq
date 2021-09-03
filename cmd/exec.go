package cmd

import (
	"fmt"

	"github.com/jlucasnsilva/httpreq"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec <filename> [...request]",
	Short: "Executes the request inside a request group YAML file",
	Long: `Executes some or all requests inside a request group YAML file.
Besides the mandatory YAML file name, you can specify any number
of request IDs. If none are provided, all will be executed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		g, err := httpreq.Parse(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		g.Execute(args[1:]...)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
