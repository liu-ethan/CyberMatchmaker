/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package service

import (
	"CyberMatchmaker/config"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// 抽出一个复用的 HTTP Client
var httpClient = &http.Client{Timeout: 30 * time.Second}

// GenerateFortune 调用文本大模型进行算命推演
func GenerateFortune(prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model": config.AppConfig.LLM.Model, // 从配置读取模型：gpt-4o-mini
		"messages": []map[string]string{
			{"role": "system", "content": "你是一个精通中国传统八字和五行的命理大师。"},
			{"role": "user", "content": prompt},
		},
	}
	payload, _ := json.Marshal(reqBody)

	// 【关键修改】使用你 yaml 里的 base_url 动态拼接请求地址
	reqURL := config.AppConfig.LLM.BaseURL + "/chat/completions"
	req, _ := http.NewRequest("POST", reqURL, bytes.NewReader(payload))

	// 使用你 yaml 里的 api_key
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.LLM.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// 增加简单的容错，防止中转站返回报错信息（比如余额不足）导致强转 panic
	if result["choices"] == nil {
		return "", errors.New("LLM 请求失败，返回信息: " + string(body))
	}

	choices := result["choices"].([]interface{})
	firstChoice := choices[0].(map[string]interface{})
	message := firstChoice["message"].(map[string]interface{})

	return message["content"].(string), nil
}

// GenerateEmbedding 调用 Embedding 模型生成 1536 维向量
func GenerateEmbedding(text string) ([]float32, error) {
	reqBody := map[string]interface{}{
		"model": config.AppConfig.LLM.EmbeddingModel, // 从配置读取模型：text-embedding-3-small
		"input": text,
	}
	payload, _ := json.Marshal(reqBody)

	// 【关键修改】使用你 yaml 里的 base_url 动态拼接请求地址
	reqURL := config.AppConfig.LLM.BaseURL + "/embeddings"
	req, _ := http.NewRequest("POST", reqURL, bytes.NewReader(payload))

	req.Header.Set("Authorization", "Bearer "+config.AppConfig.LLM.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, errors.New("empty embedding result")
	}

	return result.Data[0].Embedding, nil
}
