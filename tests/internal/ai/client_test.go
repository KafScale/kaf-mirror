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


package ai_test

import (
	"context"
	"kaf-mirror/internal/ai"
	"kaf-mirror/internal/config"
	"kaf-mirror/tests/mocks"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestAIClient(t *testing.T) {
	mock := &mocks.MockOpenAIClient{
		CreateChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
			return openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: "Test response",
						},
					},
				},
			}, nil
		},
	}

	client := &ai.Client{
		Client: mock,
		Cfg:    config.AIConfig{Model: "test-model"},
	}

	t.Run("GetAnomalyDetection", func(t *testing.T) {
		resp, err := client.GetAnomalyDetection(context.Background(), "some metrics")
		assert.NoError(t, err)
		assert.Equal(t, "Test response", resp)
	})

	t.Run("GetPerformanceRecommendation", func(t *testing.T) {
		resp, err := client.GetPerformanceRecommendation(context.Background(), "some metrics")
		assert.NoError(t, err)
		assert.Equal(t, "Test response", resp)
	})

	t.Run("ExplainEvent", func(t *testing.T) {
		resp, err := client.ExplainEvent(context.Background(), "some event")
		assert.NoError(t, err)
		assert.Equal(t, "Test response", resp)
	})
}
