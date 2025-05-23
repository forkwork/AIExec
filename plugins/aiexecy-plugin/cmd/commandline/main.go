package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCommand = &cobra.Command{
		Use:   "aiexec",
		Short: "Aiexec",
		Long:  "Aiexec is a cli tool to help you develop your Aiexec projects.",
	}

	pluginCommand = &cobra.Command{
		Use:   "plugin",
		Short: "Plugin",
		Long:  "Plugin related commands",
	}

	bundleCommand = &cobra.Command{
		Use:   "bundle",
		Short: "Bundle",
		Long:  "Bundle related commands",
	}

	signatureCommand = &cobra.Command{
		Use:   "signature",
		Short: "Signature",
		Long:  "Signature related commands",
	}

	versionCommand = &cobra.Command{
		Use:   "version",
		Short: "Version",
		Long:  "Show the version of aiexec cli",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(VersionX)
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCommand.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aiexec.yaml)")
	rootCommand.AddCommand(pluginCommand)
	rootCommand.AddCommand(bundleCommand)
	rootCommand.AddCommand(signatureCommand)
	rootCommand.AddCommand(versionCommand)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".aiexec" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".aiexec")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	rootCommand.Execute()
}
