package main

import (
	"github.com/Calmantara/lis-backend/internal/adaptors/http"
	"github.com/spf13/cobra"
)

const (
	HTTP_COMMAND = "This command will start the lis HTTP API server to handle requests."
)

var httpCommand = &cobra.Command{
	Use:   "http",
	Short: "Run lis HTTP API server",
	Long:  HTTP_COMMAND,
	Run:   http.RunEcho,
}
