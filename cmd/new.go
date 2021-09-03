package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jlucasnsilva/httpreq"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [filenames]",
	Short: "Creates yaml files with the basic template",
	Long: `Create new yaml files with the given names if they don't exist.
The new files will contain the basic request group template.`,
	Run: newCmdRun,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func newCmdRun(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("At least one file name is required")
		return
	}

	for _, fileName := range args {
		fname := fileName
		if !strings.HasSuffix(fname, ".yml") {
			fname += ".yml"
		}

		if _, err := os.Stat(fname); err == nil {
			fmt.Printf("file '%v' already exists\n", fname)
			continue
		}

		if err := httpreq.CreateFile(fname); err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("created '%v'\n", fname)
	}
}
