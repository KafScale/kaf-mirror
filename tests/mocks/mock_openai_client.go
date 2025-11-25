// kaf-mirror - A high-performance Kafka replication tool with AI-powered operational intelligence.
// Copyright (C) 2025 Scalytics
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.


package mocks

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// MockOpenAIClient is a mock implementation of the OpenAIClient interface.
type MockOpenAIClient struct {
	CreateChatCompletionFunc func(context.Context, openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	if m.CreateChatCompletionFunc != nil {
		return m.CreateChatCompletionFunc(ctx, req)
	}
	return openai.ChatCompletionResponse{}, nil
}
