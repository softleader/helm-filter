package main

import (
	"errors"
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"path/filepath"
	"path"
)

const desc = `
Filter out template files
	$ helm filter mychart"
`

func main() {
	filterCmd := filterCmd{}

	cmd := &cobra.Command{
		Use:   "helm filter [flags] CHART",
		Short: fmt.Sprintf("filter out template files"),
		Long:  desc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("chart is required")
			}
			// verify chart path exists
			if _, err := os.Stat(args[0]); err == nil {
				if filterCmd.chartPath, err = filepath.Abs(args[0]); err != nil {
					return err
				}
			} else {
				return err
			}
			// verify filter file exists
			if !filepath.IsAbs(filterCmd.valuesFile) {
				filterCmd.valuesFile = path.Join(filterCmd.chartPath, filterCmd.valuesFile)
			}
			_, err := os.Stat(filterCmd.valuesFile)
			if os.IsNotExist(err) {
				return fmt.Errorf("values '%s' does not exist", filterCmd.valuesFile)
			}
			return filterCmd.run()
		},
	}
	f := cmd.Flags()
	f.StringVarP(&filterCmd.outputDir, "output-dir", "o", "", "copy all files to output-dir and filter there instead filter in chart path")
	f.StringVarP(&filterCmd.valuesFile, "values", "f", "values.yaml", "specify values in a YAML file to filter")
	f.BoolVarP(&filterCmd.overwriteValues, "overwrite-values", "", false, "overwrite values after filtered out")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
