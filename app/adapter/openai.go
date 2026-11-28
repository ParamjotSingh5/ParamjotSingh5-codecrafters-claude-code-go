package adapter

import (
	"github.com/codecrafters-io/claude-code-starter-go/app/tools"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
)

func ToOpenAITools(toolDefs []tools.Tool) []openai.ChatCompletionToolUnionParam {
	result := make([]openai.ChatCompletionToolUnionParam, 0, len(toolDefs))
	for _, def := range toolDefs {
		result = append(result, openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
			Name:        def.Name(),
			Description: param.NewOpt(def.Description()),
			Parameters:  def.ParametersSchema(),
		}))
	}
	return result
}
