package cmd

import (
	"fmt"
	"github.com/nginxinc/kubectl-kic-supportpkg/pkg/data_collector"
	"github.com/nginxinc/kubectl-kic-supportpkg/pkg/jobs"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {

	var namespaces []string

	var rootCmd = &cobra.Command{
		Use:   "kic-supportpkg",
		Short: "kic-supportpkg - a tool to create Ingress Controller diagnostics package",
		Long:  `kic-supportpkg - a tool to create Ingress Controller diagnostics package`,
		Run: func(cmd *cobra.Command, args []string) {

			collector, err := data_collector.NewDataCollector(namespaces...)
			if err != nil {
				fmt.Println(fmt.Errorf("unable to start data collector: %s", err))
				os.Exit(1)
			}

			for _, job := range jobs.JobList() {
				fmt.Printf("Running job %s...", job.Name)
				err = job.Collect(collector)
				if err != nil {
					fmt.Printf("Error: %s", err)
				} else {
					fmt.Print(" OK\n")
				}
			}

			err = collector.WrapUp()
			if err != nil {
				fmt.Println(fmt.Errorf("error when wrapping up: %s", err))
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringSliceVarP(&namespaces, "namespace", "n", []string{}, "comma-separated list of namespaces to collect information from")
	rootCmd.SetUsageTemplate("Usage: \n kic-supportpkg -n namespace1 -n namespace2 ...")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
