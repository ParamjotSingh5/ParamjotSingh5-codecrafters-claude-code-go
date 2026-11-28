package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/codecrafters-io/claude-code-starter-go/app/adapter"
	"github.com/codecrafters-io/claude-code-starter-go/app/tools"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("OPENROUTER_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
	}

	toolManager := tools.NewManager()
	toolManager.Register(&tools.ReadTool{})

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(prompt),
	}

	for {
		params := openai.ChatCompletionNewParams{
			Model:    "anthropic/claude-haiku-4.5",
			Messages: messages,
		}

		if toolsDefs := toolManager.List(); len(toolsDefs) > 0 {
			params.Tools = adapter.ToOpenAITools(toolsDefs)
		}

		resp, err := client.Chat.Completions.New(context.Background(), params)

		if err != nil {
			os.Exit(1)
		}

		if len(resp.Choices) == 0 {
			panic("No choices in response")
		}

		choice := resp.Choices[0]

		messages = append(messages, choice.Message.ToParam())

		// You can use print statements as follows for debugging, they'll be visible when running tests.
		fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

		fmt.Print(resp.Choices[0].Message.Content)

		if len(choice.Message.ToolCalls) == 0 {
			fmt.Println("No tool call requested, exiting agent-loop")
			break
		}

		for _, toolCall := range choice.Message.ToolCalls {
			toolName := toolCall.Function.Name
			toolInput := toolCall.Function.Arguments

			result, err := toolManager.ExecuteCall(toolName, toolInput)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}

			messages = append(messages, openai.ToolMessage(result, toolCall.ID))
		}

	}

	os.Exit(0)
}
