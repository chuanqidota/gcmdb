package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gcmdb/config"
	"gcmdb/pkg/logger"
	"io"
	"net/http"
	"time"
)

type llmClient struct{}

var LLMClient = new(llmClient)

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest OpenAI 兼容请求
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

// ChatCompletionResponse OpenAI 兼容响应
type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Chat 发送聊天请求到大模型 API
func (lc *llmClient) Chat(messages []ChatMessage) (string, int, error) {
	cfg := config.Conf.LLM
	if cfg.APIURL == "" || cfg.APIKey == "" {
		return "", 0, fmt.Errorf("大模型服务未配置，请在 config.yaml 中配置 llm.api_url 和 llm.api_key")
	}

	maxTokens := cfg.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	reqBody := ChatCompletionRequest{
		Model:       cfg.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: cfg.Temperature,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, fmt.Errorf("序列化请求失败: %s", err.Error())
	}

	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("POST", cfg.APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, fmt.Errorf("创建请求失败: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("LLM API 调用失败: %s", err.Error()))
		return "", 0, fmt.Errorf("大模型服务暂不可用，请稍后重试")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("读取响应失败: %s", err.Error())
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", 0, fmt.Errorf("解析响应失败: %s", err.Error())
	}

	if chatResp.Error != nil {
		return "", 0, fmt.Errorf("大模型返回错误: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", 0, fmt.Errorf("大模型未返回有效回复")
	}

	return chatResp.Choices[0].Message.Content, chatResp.Usage.TotalTokens, nil
}
