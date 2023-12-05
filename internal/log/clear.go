package log

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fanniva/resy-cli/internal/utils/paths"
)

func Clear() error {
	p, err := paths.GetAppPaths()
	if err != nil {
		return err
	}
	confirm := false
	err = survey.AskOne(&survey.Confirm{Message: "Clear all log files?"}, &confirm)
	if err != nil {
		return err
	}

	if confirm {
		return os.RemoveAll(p.LogPath)
	}

	return nil
}
