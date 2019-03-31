package main

import (
	"fmt"
	"github.com/nicored/forge_tools/properties"

	"github.com/spf13/cobra"
)

func main() {
	cmd := newCmd()
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "",
	}

	addPropertiesCmd(cmd)

	return cmd
}

func addPropertiesCmd(cmd *cobra.Command) {
	propCmd := &cobra.Command{
		Use:   "properties",
		Short: "Operate on model properties",
	}

	toJsonCmd := &cobra.Command{
		Use:     "json [objects_dir_path]",
		Short:   "Extracts properties from objects files and prints them to json format",
		Example: "properties json ~/Documents/derivatives/example/properties",
		Run: func(jcmd *cobra.Command, args []string) {
			// defaults to current directory
			path := "."

			if len(args) > 0 {
				path = args[0]
			}

			props := properties.NewProperties(path)

			bJson, err := props.ExportJson(true)
			if err != nil {
				fmt.Println("error", err)
			}

			fmt.Println(string(bJson))
		},
	}

	propCmd.AddCommand(toJsonCmd)
	cmd.AddCommand(propCmd)
}
