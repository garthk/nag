package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type renderer func(cmd *cobra.Command, dirname string) error

var renderers map[string]renderer

var bragManPageSection string
var bragOutputDirectory string
var bragChosenFormat string

func init() {
	renderers = make(map[string]renderer)
	renderers["md"] = genMarkdownTree
	renderers["man"] = genManTree

	bragCmd.Flags().StringVarP(&bragOutputDirectory, "output", "o", "", "output directory (default temporary)")
	bragCmd.Flags().StringVarP(&bragChosenFormat, "format", "f", "md", "output format (md|man)")
	bragCmd.Flags().StringVarP(&bragManPageSection, "section", "s", "8", "man page section")
	RootCmd.AddCommand(bragCmd)
}

var bragCmd = &cobra.Command{
	Use:    "brag",
	Short:  "Generate documentation for command line help",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		renderer := renderers[bragChosenFormat]
		if renderer == nil {
			return errors.New(fmt.Sprintf("unknown format: %s", bragChosenFormat))
		}

		dirname, err := checkOrCreate(bragOutputDirectory, RootCmd.CommandPath())
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		fmt.Printf("rendering %s to: %s\n", bragChosenFormat, dirname)
		return renderer(RootCmd, dirname)
	},
}

func checkOrCreate(dirname, prefix string) (string, error) {
	if dirname == "" {
		return ioutil.TempDir("", prefix)
	} else {
		info, err := os.Stat(dirname)
		if err != nil {
			return dirname, err
		} else if !info.IsDir() {
			return dirname, errors.New(fmt.Sprintf("not a directory: %s", dirname))
		} else {
			return dirname, nil
		}
	}
}

func genMarkdownTree(cmd *cobra.Command, dirname string) error {
	return doc.GenMarkdownTree(cmd, dirname)
}

func genManTree(cmd *cobra.Command, dirname string) error {
	header := &doc.GenManHeader{
		Title:   "",  // auto-generate
		Section: "8", // system administration tools and procedures
		Source:  getPackagePath(),
	}
	return doc.GenManTree(cmd, header, dirname)
}

type pkgSensor struct{}

func getPackagePath() string {
	cmdPkgPath := reflect.TypeOf(pkgSensor{}).PkgPath()

	// now chomp off the /cmd/ part
	lastSlash := strings.LastIndex(cmdPkgPath, "/")
	if lastSlash < 0 {
		return cmdPkgPath
	} else {
		return cmdPkgPath[:lastSlash]
	}
}
