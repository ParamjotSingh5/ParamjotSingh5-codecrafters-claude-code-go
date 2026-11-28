package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

type ReadInput struct {
	FilePath string `json:"file_path"`
}
type ReadTool struct{}

func (r *ReadTool) Name() string {
	return "read_file"
}

func (r *ReadTool) Description() string {
	return "Read and return the contents of a file"
}

func (r *ReadTool) ParametersSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "The path to the file to read",
			},
		},
		"required": []string{"file_path"},
	}
}

func (r *ReadTool) Execute(input json.RawMessage) (string, error) {
	var params ReadInput

	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	var result map[string]any
	err := json.Unmarshal([]byte(params.FilePath), &result)
	if err != nil {
		return "", err
	}

	filePathArg := result["file_path"].(string)

	data, err := os.ReadFile(filePathArg)

	if err != nil {

		return "", err
	}

	return string(data), nil
}
