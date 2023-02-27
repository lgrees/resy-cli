package paths

import (
	"os"
	"path"

	"github.com/lgrees/resy-cli/constants"
)

type paths struct {
	AppPath        string
	LogPath        string
	ConfigFilePath string
}

func GetAppPaths() (*paths, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appPath := path.Join(configPath, constants.AppDirName)
	configFilePath := path.Join(appPath, constants.AppAuthCfgFile)
	logPath := path.Join(appPath, constants.AppLogDir)

	return &paths{
		AppPath:        appPath,
		ConfigFilePath: configFilePath,
		LogPath:        logPath,
	}, nil
}
