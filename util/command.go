package util

import (
	"bufio"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log/slog"
	"os/exec"
	"strings"
)

func convertGBKToUTF8(r io.Reader) io.Reader {
	return transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
}

func ExecuteCommand(command string, args []string, printErrorOnly bool) {
	slog.Info("Executing command: " + command + " " + strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("Error creating stderr pipe: " + err.Error())
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Error creating stdout pipe: " + err.Error())
		return
	}
	if err := cmd.Start(); err != nil {
		slog.Error("Error starting command: " + err.Error())
		return
	}
	go func() {
		scanner := bufio.NewScanner(convertGBKToUTF8(stderr))
		for scanner.Scan() {
			slog.Debug(scanner.Text())
		}
	}()
	go func() {
		scanner := bufio.NewScanner(convertGBKToUTF8(stdout))
		for scanner.Scan() {
			text := scanner.Text()
			if !printErrorOnly || strings.Contains(text, "ERROR") {
				slog.Debug(text)
			}
		}
	}()
	if err := cmd.Wait(); err != nil {
		slog.Error("Command finished with error: " + err.Error())
	}
}
