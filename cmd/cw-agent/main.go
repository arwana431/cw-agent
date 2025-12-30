// Package main provides the entry point for the CertWatch Agent.
package main

import (
	"os"

	"github.com/certwatch-app/cw-agent/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
