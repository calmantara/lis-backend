package main

import (
	"github.com/Calmantara/lis-backend/internal/adaptors/jobs"
	"github.com/spf13/cobra"
)

const (
	JOB_COMMAND       = "This command will show how to run lis job related commands (e.g. Migration, Scheduler)"
	MIGRATION_COMMAND = "This command will run the database migration to ensure the database schema is up to date."
)

var jobCommand = &cobra.Command{
	Use:   "job",
	Short: "Run lis job related commands",
	Long:  JOB_COMMAND,
}

var migrationCommand = &cobra.Command{
	Use:   "migration",
	Short: "Run database migration",
	Long:  MIGRATION_COMMAND,
	Run:   jobs.Migrate,
}
