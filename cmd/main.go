package main

import (
	"github.com/aquasecurity/trivy/pkg/utils/fsutils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"trivy-plugin-vulners/internal"
)

func main() {
	Execute()
}
func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var (
	cacheDir string
	apiKey   string
)

var rootCmd = &cobra.Command{
	Use:          "trivy-plugin-vulners-db",
	Short:        "trivy-plugin-vulners-db",
	Long:         "trivy-plugin-vulners-db",
	SilenceUsage: true,
	Version:      "0.1",
	Args:         cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(os.Stdout, "info"); err != nil {
			return err
		}
		if len(apiKey) == 0 {
			log.Fatalf("Missing api key")
		}

		if len(cacheDir) == 0 {
			cacheDir = fsutils.CacheDir()
		}
		internal.Download(cacheDir, apiKey)
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&cacheDir, "cache-dir", "", "", "cache dir")
	rootCmd.Flags().StringVarP(&apiKey, "api-key", "", "", "vulners api key")
}

func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}
