package main

import (
	"errors"
	"fmt"

	"github.com/raymondji/commitstack/sampleusage"
	"github.com/spf13/cobra"
)

var learnPartFlag int

func init() {
	learnCmd.Flags().IntVar(&learnPartFlag, "part", 1, "Which part of the tutorial to continue from")
}

var learnCmd = &cobra.Command{
	Use:   "learn",
	Short: "Prints sample commands to learn how to use commitstack",
	RunE: func(cmd *cobra.Command, args []string) error {
		deps, err := initDeps()
		if err != nil {
			return err
		}
		samples := sampleusage.New(deps.theme, deps.repoCfg.DefaultBranch, deps.git, deps.host)

		switch learnPartFlag {
		case 1:
			if err := samples.Cleanup(); err != nil {
				fmt.Println("ERROR failed to cleanup sample", err.Error())
				return nil
			}
			if err := samples.Part1().Execute(); err != nil {
				fmt.Println("ERROR failed to execute sample", err.Error())
			}
		case 2:
			fmt.Println("TODO")
		default:
			return errors.New("invalid tutorial part")
		}

		return nil
	},
}
