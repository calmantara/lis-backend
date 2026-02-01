package main

import (
	"fmt"
	"os"

	"github.com/Calmantara/lis-backend/internal/helpers/errors"
	"github.com/spf13/cobra"
)

const (
	version = "0.0.1"
	cliName = "lis-backend"
)

var rootCmd = &cobra.Command{
	Use:     cliName,
	Version: version,
	Short:   "Core CLI to run laboratory information service",
	Long:    "This command will show how to run lis API / worker / job",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, errors.ERROR_MAIN_COMMAND.Error()))
		os.Exit(1)
	}
}

func init() {
	// init all commands
	jobCommand.AddCommand(migrationCommand)

	// manage all commands
	rootCmd.AddCommand(httpCommand, jobCommand)
}
