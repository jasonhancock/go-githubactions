package githubactions

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
)

// AddMask masks a value from logs.
// https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#masking-a-value-in-log
func AddMask(value string) error {
	return FAddMask(os.Stdout, value)
}

// FAddMask masks a value from logs. Writes the command to the specified
// io.Writer. Useful if you want to send your workflow commands to
// os.Stderr.
func FAddMask(w io.Writer, value string) error {
	_, err := fmt.Fprintf(w, "::add-mask::%s\n", value)
	return err

}

// SetOutput sets an output variable. Appends the value to the $GITHUB_OUTPUT
// file.
func SetOutput(name, value string) error {
	file := os.Getenv("GITHUB_OUTPUT")
	if file == "" {
		return errors.New("GITHUB_OUTPUT env variable not specified")
	}

	return appendToFile(file, fmt.Sprintf("%s=%s\n", name, value))
}

func appendToFile(file, content string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = f.WriteString(content)
	closeErr := f.Close()
	if err != nil {
		return err
	}

	return closeErr
}

// SetSummary appends markdown content to the $GITHUB_STEP_SUMMARY file.
func SetSummary(markdownContent string) error {
	file := os.Getenv("GITHUB_STEP_SUMMARY")
	if file == "" {
		return errors.New("GITHUB_STEP_SUMMARY env variable not specified")
	}

	if !strings.HasSuffix(markdownContent, "\n") {
		markdownContent += "\n"
	}

	return appendToFile(file, markdownContent)
}

// LogDebugf logs a debug message.
func LogDebugf(format string, a ...any) error {
	return FLogDebugf(os.Stdout, format, a...)
}

// FLogDebugf logs a debug message to the specified io.Writer.
func FLogDebugf(w io.Writer, format string, a ...any) error {
	_, err := fmt.Fprintf(w, "::debug::"+format+"\n", a...)
	return err
}

// LogWarnf logs a warning message.
func LogWarnf(format string, a ...any) error {
	return FLogWarnf(os.Stdout, format, a...)
}

// FLogWarnf logs a warning message to the specified io.Writer.
func FLogWarnf(w io.Writer, format string, a ...any) error {
	_, err := fmt.Fprintf(w, "::warning::"+format+"\n", a...)
	return err
}

// LogErrorf logs an error message.
func LogErrorf(format string, a ...any) error {
	return FLogErrorf(os.Stdout, format, a...)
}

// FLogErrorf logs an error message to the specified io.Writer.
func FLogErrorf(w io.Writer, format string, a ...any) error {
	_, err := fmt.Fprintf(w, "::error::"+format+"\n", a...)
	return err
}

// GetRunURL attemtps to construct a URL to the action run from environment
// variables exposed during workflow execution.
func GetRunURL() (string, error) {
	u, err := url.Parse(os.Getenv("GITHUB_SERVER_URL"))
	if err != nil {
		return "", fmt.Errorf("parsing server url: %w", err)
	}

	u.Path = path.Join(os.Getenv("GITHUB_REPOSITORY"), "actions", "runs", os.Getenv("GITHUB_RUN_ID"))

	return u.String(), nil
}
