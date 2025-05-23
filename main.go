package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// Import all the tools implemented here.
	_ "github.com/acuvity/mcp-server-microsoft-graph/api/applications"
	_ "github.com/acuvity/mcp-server-microsoft-graph/api/sites"
	_ "github.com/acuvity/mcp-server-microsoft-graph/api/users"
	"github.com/acuvity/mcp-server-microsoft-graph/cmd/cli"
	"github.com/acuvity/mcp-server-microsoft-graph/mcp"
)

func main() {

	name := "mcp-server-microsoft-graph"
	description := "Microsoft MCP Command Line Tool"
	version := "1.0.0"

	cobra.OnInitialize(func() {
		viper.SetEnvPrefix(name)
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	})

	var rootCmd = &cobra.Command{
		Use:   name,
		Short: description,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version and exit.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	var cliCommand = &cobra.Command{
		Use:   "cli",
		Short: "Run CLI.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		Run: func(cmd *cobra.Command, args []string) {
			cli.Run(cmd, args)
		},
	}

	rootCmd.AddCommand(
		versionCmd,
		cliCommand,
	)

	rootCmd.PersistentFlags().String("tenant-id", "", "Microsoft Tenant ID")
	rootCmd.PersistentFlags().String("client-id", "", "Microsoft Client ID")
	rootCmd.PersistentFlags().String("client-secret", "", "Microsoft Client Secret")
	rootCmd.PersistentFlags().String("transport", "sse", "MCP transport type (stdio or sse)")
	rootCmd.PersistentFlags().String("service-name", "localhost", "Microsoft Service Name")

	viper.SetConfigName("config") // name of the file (without extension)
	viper.SetConfigType("yaml")   // or viper.SetConfigType("json") if it's json
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	// Read in the config
	_ = viper.ReadInConfig()

	rootCmd.RunE = mcp.Run
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
