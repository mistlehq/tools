package textinput

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Read normalizes a mutually exclusive `--value`/`--file` pair and supports
// `-` as stdin for file-backed input.
func Read(stdin io.Reader, valueFlagName string, value string, fileFlagName string, filePath string) (string, error) {
	if value != "" && filePath != "" {
		return "", fmt.Errorf("--%s and --%s are mutually exclusive", valueFlagName, fileFlagName)
	}

	if value == "" && filePath == "" {
		return "", fmt.Errorf("exactly one of --%s or --%s is required", valueFlagName, fileFlagName)
	}

	if value != "" {
		if strings.TrimSpace(value) == "" {
			return "", fmt.Errorf("--%s must not be empty", valueFlagName)
		}

		return value, nil
	}

	var body []byte
	var err error

	switch filePath {
	case "-":
		body, err = io.ReadAll(stdin)
	default:
		body, err = os.ReadFile(filePath)
	}

	if err != nil {
		return "", err
	}

	text := string(body)
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("--%s must not be empty", fileFlagName)
	}

	return text, nil
}
