/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package middleware

import (
	"CyberMatchmaker/config"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type LLMService struct {
	Model llms.Model
}

var LLM *LLMService

// NewLLMService 工厂函数：在这里把臃肿的配置写好
func NewLLMService() error {
	instance, err := openai.New(
		openai.WithToken(config.AppConfig.LLM.APIKey),
		openai.WithBaseURL(config.AppConfig.LLM.BaseURL),
		openai.WithModel(config.AppConfig.LLM.Model),
	)
	if err != nil {
		return nil
	}
	LLM = &LLMService{Model: instance}
	return nil
}

// CallAI 统一的调用入口
func (s *LLMService) CallAI(ctx context.Context, sysPrompt, userPrompt string) (string, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, sysPrompt),
		llms.TextParts(llms.ChatMessageTypeGeneric, userPrompt),
	}
	resp, err := s.Model.GenerateContent(ctx, content)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI 返回结果为空")
	}
	return resp.Choices[0].Content, nil
}
