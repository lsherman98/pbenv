package routes

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func renderEnvPageHandler(e *core.RequestEvent) error {
	html, err := template.NewRegistry().LoadFiles(
		"views/layout.html",
		"views/env.html",
	).Render(nil)
	if err != nil {
		return e.BadRequestError("failed to load html", err)
	}

	return e.HTML(http.StatusOK, html)
}

func addEnvHandler(e *core.RequestEvent) error {
	var newVar EnvVar
	if err := e.BindBody(&newVar); err != nil {
		return e.BadRequestError("invalid request body", err)
	}

	if newVar.Key == "" {
		return e.BadRequestError("key is required", nil)
	}

	envVars, err := readEnvFile()
	if err != nil {
		return e.BadRequestError("failed to read .env file", err)
	}

	for _, env := range envVars {
		if env.Key == newVar.Key {
			return e.BadRequestError("key already exists", nil)
		}
	}

	lines, err := readLines()
	if err != nil {
		return e.BadRequestError("failed to read .env file", err)
	}

	lines = append(lines, fmt.Sprintf("%s=%s", newVar.Key, newVar.Value))
	if err := writeLines(lines); err != nil {
		return e.BadRequestError("failed to write to .env file", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Environment variable added successfully",
	})
}

func updateEnvHandler(e *core.RequestEvent) error {
	key := e.Request.PathValue("key")
	if key == "" {
		return e.BadRequestError("key is required", nil)
	}

	var updateVar EnvVar
	if err := e.BindBody(&updateVar); err != nil {
		return e.BadRequestError("invalid request body", err)
	}

	lines, err := readLines()
	if err != nil {
		return e.BadRequestError("failed to read .env file", err)
	}

	envVars, _ := readEnvFile()
	found := false
	for _, env := range envVars {
		if env.Key == key {
			found = true
			break
		}
	}

	if !found {
		return e.NotFoundError("environment variable not found", nil)
	}

	lines, err = updateEnv(lines, key, updateVar.Value)
	if err != nil {
		return e.BadRequestError("failed to update environment variable", err)
	}

	if err := writeLines(lines); err != nil {
		return e.BadRequestError("failed to write to .env file", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "environment variable updated successfully",
	})
}

func deleteEnvHandler(e *core.RequestEvent) error {
	key := e.Request.PathValue("key")
	if key == "" {
		return e.BadRequestError("key is required", nil)
	}

	lines, err := readLines()
	if err != nil {
		return e.BadRequestError("failed to read .env file", err)
	}

	envVars, _ := readEnvFile()
	found := false
	for _, env := range envVars {
		if env.Key == key {
			found = true
			break
		}
	}

	if !found {
		return e.NotFoundError("environment variable not found", err)
	}

	lines, err = deleteEnv(lines, key)
	if err != nil {
		return e.BadRequestError("failed to delete environment variable", err)
	}

	if err := writeLines(lines); err != nil {
		return e.BadRequestError("failed to write to .env file", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "environment variable deleted successfully",
	})
}

func getEnvsHandler(e *core.RequestEvent) error {
	envVars, err := readEnvFile()
	if err != nil {
		return e.BadRequestError("failed to read environment variables", err)
	}

	return e.JSON(http.StatusOK, envVars)
}

func readEnvFile() ([]EnvVar, error) {
	envVars := make([]EnvVar, 0)

	lines, err := readLines()
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, `"'`)
			envVars = append(envVars, EnvVar{Key: key, Value: value})
		}
	}

	return envVars, nil
}

func readLines() ([]string, error) {
	file, err := os.Open(".env")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func writeLines(lines []string) error {
	file, err := os.Create(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		_, err := fmt.Fprintln(file, line)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateEnv(lines []string, key, value string) ([]string, error) {
	success := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		comment := ""
		lineWithoutComment := line
		if idx := strings.Index(line, "#"); idx != -1 && idx > 0 {
			comment = line[idx:]
			lineWithoutComment = strings.TrimSpace(line[:idx])
		}

		parts := strings.SplitN(lineWithoutComment, "=", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			if comment != "" {
				lines[i] = fmt.Sprintf("%s=%s %s", key, value, comment)
			} else {
				lines[i] = fmt.Sprintf("%s=%s", key, value)
			}
			success = true
			break
		}
	}

	if !success {
		return lines, fmt.Errorf("environment variable not found")
	}

	return lines, nil
}

func deleteEnv(lines []string, key string) ([]string, error) {
	newLines := make([]string, 0)
	success := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			newLines = append(newLines, line)
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			success = true
			continue
		}

		newLines = append(newLines, line)
	}

	if !success {
		return lines, fmt.Errorf("environment variable not found")
	}

	return newLines, nil
}
