/**
 * @author 刘潇翰
 * @since 2026/3/2
 */
package utils

import (
	"fmt"
	"strings"

	"github.com/goccy/go-json"
)

// StringtoClass 负责将 AI 的 String 结果直接“缝合”到现有的 Record 中
func StringtoClass(aiRaw string, record interface{}) error {
	// 1. 清洗：去掉 AI 常见的 Markdown 标签
	cleanJSON := strings.TrimSpace(aiRaw)
	cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")
	cleanJSON = strings.TrimSpace(cleanJSON)
	// 2. 填充：Unmarshal 到指针，实现增量更新
	// 注意：record 必须是一个指针
	return json.Unmarshal([]byte(cleanJSON), record)
}

// CleanMarkdown 接收字符串指针，原地修改内容并返回提取状态
func CleanMarkdown(raw *string) error {
	if raw == nil || *raw == "" {
		return fmt.Errorf("输入字符串为空")
	}
	// 1. 寻找 JSON 对象的边界
	s := *raw
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	// 2. 校验边界合法性
	if start == -1 || end == -1 {
		return fmt.Errorf("未在内容中找到 JSON 数据块 (缺少 { 或 })")
	}
	if start >= end {
		return fmt.Errorf("JSON 数据块格式异常 (反括号在正括号之前)")
	}
	// 3. 截取并原地更新指针指向的内容
	// strings.TrimSpace 确保不会留下换行符
	*raw = strings.TrimSpace(s[start : end+1])
	return nil
}
