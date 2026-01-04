package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/chat"
)

type Tool struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
	Handler     func(args map[string]interface{}) (string, error)
}

type ToolRegistry map[string]Tool

func NewToolRegistry() ToolRegistry {
	registry := make(ToolRegistry)

	registry["read_file"] = Tool{
		Name: "read_file",
		Description: "Reads the complete contents of a file at the specified path. " +
			"Use this when you need to examine file contents. " +
			"Input should be an absolute or relative file path.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The file path to read (absolute or relative)",
				},
			},
			"required": []string{"path"},
		},
		Handler: readFileHandler,
	}

	registry["list_directory"] = Tool{
		Name: "list_directory",
		Description: "Lists all files and directories at the specified path. " +
			"Use this to explore directory structure or find files. " +
			"Returns a list of names with directories marked by trailing /. " +
			"If no path is provided, lists the current directory.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The directory path to list (defaults to current directory if not provided)",
				},
			},
		},
		Handler: listDirectoryHandler,
	}

	registry["write_file"] = Tool{
		Name: "write_file",
		Description: "Writes content to a file at the specified path. " +
			"Creates the file if it doesn't exist, overwrites if it does. " +
			"Use this to create or modify files based on user requests. " +
			"WARNING: This operation will overwrite existing files without confirmation.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The file path to write to",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "The content to write to the file",
				},
			},
			"required": []string{"path", "content"},
		},
		Handler: writeFileHandler,
	}

	return registry
}

func (tr ToolRegistry) ToSDKTools() []chat.Tool {
	tools := make([]chat.Tool, 0, len(tr))

	for name, tool := range tr {
		sdkTool := chat.NewFunctionTool(
			name,
			tool.Description,
			tool.Parameters,
		)
		tools = append(tools, sdkTool)
	}

	return tools
}

func readFileHandler(args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("path parameter must be a string")
	}

	if path == "" {
		return "", fmt.Errorf("path parameter cannot be empty")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not found: %s", path)
		}
		if os.IsPermission(err) {
			return "", fmt.Errorf("permission denied: %s", path)
		}
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return string(content), nil
}

func listDirectoryHandler(args map[string]interface{}) (string, error) {
	path := "."
	if p, ok := args["path"].(string); ok && p != "" {
		path = p
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("directory not found: %s", path)
		}
		if os.IsPermission(err) {
			return "", fmt.Errorf("permission denied: %s", path)
		}
		return "", fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	var result strings.Builder
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		result.WriteString(name)
		result.WriteString("\n")
	}

	output := result.String()
	if output == "" {
		return "(empty directory)", nil
	}

	return output, nil
}

func writeFileHandler(args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("path parameter must be a string")
	}
	if path == "" {
		return "", fmt.Errorf("path parameter cannot be empty")
	}

	content, ok := args["content"].(string)
	if !ok {
		return "", fmt.Errorf("content parameter must be a string")
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		if os.IsPermission(err) {
			return "", fmt.Errorf("permission denied: %s", path)
		}
		if os.IsExist(err) {
			return "", fmt.Errorf("cannot write to directory: %s", path)
		}
		return "", fmt.Errorf("failed to write file %s: %w", path, err)
	}

	fileSize := len(content)
	return fmt.Sprintf("Successfully wrote %d bytes to %s", fileSize, path), nil
}
