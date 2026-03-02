/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package middleware

import (
	"CyberMatchmaker/config"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type LLMService struct {
	Model    llms.Model
	Embedder embeddings.Embedder // 使用接口类型，最通用
}

var LLM *LLMService

func NewLLMService() error {
	// 1. 初始化 OpenAI 客户端
	// 注意：openai.New 返回的实例既实现了 llms.Model (聊天)，也实现了计算向量的方法
	client, err := openai.New(
		openai.WithToken(config.AppConfig.LLM.APIKey),
		openai.WithBaseURL(config.AppConfig.LLM.BaseURL),
		openai.WithModel(config.AppConfig.LLM.Model),
	)
	if err != nil {
		return err
	}
	// 2. 使用 langchaingo 提供的构造函数创建一个 Embedder
	// 它会自动封装 client，并提供高级方法如 EmbedQuery
	e, err := embeddings.NewEmbedder(client)
	if err != nil {
		return fmt.Errorf("创建 Embedder 失败: %v", err)
	}
	LLM = &LLMService{
		Model:    client,
		Embedder: e,
	}
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

// Embedding 现在直接调用封装好的高级函数
func (s *LLMService) Embedding(ctx context.Context, text string) ([]float32, error) {
	if s.Embedder == nil {
		return nil, fmt.Errorf("embedder 未初始化")
	}
	// EmbedQuery 是 langchaingo 封装好的函数：直接输入字符串，返回 []float32
	return s.Embedder.EmbedQuery(ctx, text)
}
