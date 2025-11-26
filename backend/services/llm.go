package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"backend/config"
	"backend/models"
)

// CallOllama sends a chat request to the Ollama API
func CallOllama(request models.OllamaChatRequest) (*models.OllamaChatResponse, error) {
	ollamaBaseURL := config.Config("OLLAMA_BASE_URL")
	if ollamaBaseURL == "" {
		return nil, fmt.Errorf("OLLAMA_BASE_URL is not set")
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Ollama request: %w", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/chat", ollamaBaseURL), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send Ollama request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API returned non-200 status: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.OllamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	return &response, nil
}

// CallOpenAI sends a chat request to the OpenAI API
func CallOpenAI(request models.OpenAIChatRequest) (*models.OpenAIChatResponse, error) {
	openaiBaseURL := config.Config("OPENAI_API_BASE_URL")
	openaiAPIKey := config.Config("OPENAI_API_KEY")

	if openaiBaseURL == "" {
		return nil, fmt.Errorf("OPENAI_API_BASE_URL is not set")
	}
	if openaiAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OpenAI request: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/chat/completions", openaiBaseURL), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openaiAPIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send OpenAI request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API returned non-200 status: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.OpenAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode OpenAI response: %w", err)
	}

	return &response, nil
}
