package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"route-dumper/pkg"
)

var rootCmd = &cobra.Command{
	Use:   "route-dumper [PATH]",
	Short: "",
	Long: ``,
	// Validate the args
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires path argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Searching for files containing a @Route annotation")
		searchPath := args[0]
		files := pkg.LocateFiles(searchPath)
		fmt.Println("Found ", len(files), " files")
		if len(files) < 1 {
			fmt.Println("No files containing @Route were found")
			return
		}

		fmt.Println("Parsing route annotations from files")
		pkg.ParseFiles(files)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}