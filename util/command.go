package util

import (
	"bufio"
	"os/exec"
	"strings"
)

func ExecuteCommand(command string, args []string, printErrorOnly bool) {
	LOGGER.Info("Executing command: " + command + " " + strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		LOGGER.Error("Error creating stderr pipe: " + err.Error())
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		LOGGER.Error("Error creating stdout pipe: " + err.Error())
		return
	}
	if err := cmd.Start(); err != nil {
		LOGGER.Error("Error starting command: " + err.Error())
		return
	}
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			LOGGER.Debug(scanner.Text())
		}
	}()
	if !printErrorOnly {
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				LOGGER.Debug(scanner.Text())
			}
		}()
	}
	if err := cmd.Wait(); err != nil {
		LOGGER.Error("Command finished with error: " + err.Error())
	}
}
