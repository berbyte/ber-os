// Copyright 2025 BER - ber.run
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ber

import (
	"os"

	"github.com/berbyte/ber-os/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ber",
	Short: "Ber is an AI-powered development assistant",
	Long: `Ber is a command-line tool that helps developers
with various tasks using AI capabilities powered by OpenAI.`,
}

// tuiCmd represents the terminal UI subcommand
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start the Ber TUI interface",
	Long:  `Launch the terminal user interface for interacting with Ber.`,
	Run: func(cmd *cobra.Command, args []string) {
		StartTUI()
	},
}

// webhookCmd represents the webhook server subcommand
var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Start the Ber webhook server",
	Long:  `Start the webhook server that listens for GitHub events.`,
	Run: func(cmd *cobra.Command, args []string) {
		StartWebhook()
	},
}

// Command line flags
var (
	debugMode   bool // Enable debug level logging
	verboseMode bool // Enable verbose logging output
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Initialize logger before executing commands
	logger.Init(debugMode, verboseMode)

	var exitCode int
	defer func() {
		// Ignore sync errors as they're usually harmless stderr sync issues
		_ = logger.Sync()
		os.Exit(exitCode)
	}()

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		logger.Log.Error("Command Execution Failed", zap.Error(err))
		exitCode = 1
		return
	}
}

// init registers subcommands and flags
func init() {
	// Add subcommands
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(webhookCmd)

	// Register persistent flags available to all commands
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&verboseMode, "verbose", false, "enable verbose logging")
}
