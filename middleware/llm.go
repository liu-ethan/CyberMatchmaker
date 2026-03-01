/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package middleware

import (
	"CyberMatchmaker/config"
	"CyberMatchmaker/mapper"
	"CyberMatchmaker/model"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

type AIResult struct {
	Bazi          string `json:"bazi"`
	FiveElements  string `json:"five_elements"`
	ZodiacSign    string `json:"zodiac_sign"`
	BestCity      string `json:"best_city"`
	RecentFortune string `json:"recent_fortune"`
	Description   string `json:"description"`
}

// CallAI 负责将 Record 转化为 Prompt 并调用大模型
func CallAI(record *model.FortuneRecord) (string, error) {
	systemPrompt := config.GetPrompt("fortune_task.system")
	userPromptTmpl := config.GetPrompt("fortune_task.user")

	if systemPrompt == "" || userPromptTmpl == "" {
		return "", fmt.Errorf("配置文件中未找到有效的 Prompt 模板")
	}

	// 1. 加载并渲染模板
	tmpl, err := template.New("fortune_task").Parse(userPromptTmpl)
	if err != nil {
		return "", fmt.Errorf("模板语法错误: %v", err)
	}

	var promptBody bytes.Buffer
	// 注意：BirthDate.Format 会在模板渲染时自动处理
	if err := tmpl.Execute(&promptBody, record); err != nil {
		return "", fmt.Errorf("模板渲染失败: %v", err)
	}

	// 2. 打印一下最终发给 AI 的内容（调试神器）
	zap.S().Info("最终 User Prompt: %s", promptBody.String())

	// 2. 构造 OpenAI 兼容的消息体
	requestBody := map[string]interface{}{
		"model": config.AppConfig.LLM.Model,
		"messages": []map[string]string{
			{"role": "system", "content": config.GetPrompt("fortune_task.system")},
			{"role": "user", "content": promptBody.String()},
		},
		"temperature": 0.7,
	}

	jsonData, _ := json.Marshal(requestBody)

	// 3. 发送请求
	req, _ := http.NewRequest("POST", config.AppConfig.LLM.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.LLM.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 4. 解析响应内容
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI 接口返回异常: %s", string(body))
	}

	// 定义简单的响应结构提取 content
	var aiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &aiResp); err != nil {
		return "", err
	}

	if len(aiResp.Choices) == 0 {
		return "", fmt.Errorf("AI 未返回有效内容")
	}

	return aiResp.Choices[0].Message.Content, nil
}

// 处理 FortuneRecord 的核心逻辑
func HandleFortuneLogic(record *model.FortuneRecord) error {
	// 2. 准备 Prompt 并调用大模型 (假设你已经有了调用 AI 的函数 CallAI)
	//prompt := renderTemplate(record)
	aiRawResponse, _ := CallAI(record)

	// 3. 解析 JSON
	var aiData AIResult
	// 注意：有时 AI 会返回 ```json { ... } ```，需要过滤掉这些标签
	cleanJSON := strings.TrimPrefix(aiRawResponse, "```json")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")
	cleanJSON = strings.TrimSpace(cleanJSON)

	if err := json.Unmarshal([]byte(cleanJSON), &aiData); err != nil {
		return fmt.Errorf("解析AI结果失败: %v", err)
	}

	// 4. 将 AI 结果赋值给数据库模型 (处理指针)
	record.Bazi = &aiData.Bazi
	record.FiveElements = &aiData.FiveElements
	record.ZodiacSign = &aiData.ZodiacSign
	record.BestCity = &aiData.BestCity
	record.RecentFortune = &aiData.RecentFortune
	record.Description = &aiData.Description

	// 手动更新状态
	record.Status = "completed"

	// 5. 保存回数据库
	mapper.UpdateFortuneRecord(record)

	zap.S().Info("AI 处理完成，已更新数据库记录 ID %d", record.ID)

	return nil
}
