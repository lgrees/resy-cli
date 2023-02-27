package cmd

import (
	"fmt"
	"os"

	"github.com/lgrees/resy-cli/internal/utils/paths"
	"github.com/lgrees/resy-cli/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "resy",
	Short:   "resy lets you schedule a reservation booking at your favorite restaurant at a later time",
	Long:    `resy lets you schedule a reservation booking at your favorite restaurant at a later time.`,
	Version: version.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	p, err := paths.GetAppPaths()
	cobra.CheckErr(err)

	appPath := p.AppPath
	configFilePath := p.ConfigFilePath
	logPath := p.LogPath

	if _, err = os.Stat(appPath); os.IsNotExist(err) {
		err = os.Mkdir(appPath, os.FileMode(0777))
		if err != nil {
			fmt.Println("Error creating config directory")
			return
		}
	}

	if _, err = os.Stat(logPath); os.IsNotExist(err) {
		err = os.Mkdir(logPath, os.FileMode(0777))
		if err != nil {
			fmt.Println("Error creating log directory")
			return
		}
	}

	if _, err = os.Stat(configFilePath); os.IsNotExist(err) {
		_, err = os.Create(configFilePath)
		if err != nil {
			fmt.Println("Error creating config file")
			return
		}
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFilePath)

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
