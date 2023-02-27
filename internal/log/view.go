package log

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lgrees/resy-cli/internal/utils/paths"
)

func View() error {
	p, err := paths.GetAppPaths()
	if err != nil {
		return err
	}

	files, err := os.ReadDir(p.LogPath)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errors.New("no log files exist")
	}

	fileNames := make([]string, len(files))

	for i, f := range files {
		fileNames[i] = f.Name()
	}

	err = selectAndPrint(p.LogPath, fileNames)

	if err != nil {
		return err
	}

	return nil
}

func selectAndPrint(logPath string, fileNames []string) error {
	fileName := ""
	err := survey.AskOne(&survey.Select{
		Message: "Select log file:",
		Options: fileNames,
	}, &fileName)
	if err != nil {
		return err
	}

	fileContent, err := os.ReadFile(path.Join(logPath, fileName))
	if err != nil {
		return err
	}
	fmt.Print(string(fileContent))

	confirm := false
	survey.AskOne(&survey.Confirm{Message: "View additional log files?"}, &confirm)

	if confirm {
		return selectAndPrint(logPath, fileNames)
	}

	return nil
}
