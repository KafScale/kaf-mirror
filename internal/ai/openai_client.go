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


package ai

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient is an interface that abstracts the methods we use from the go-openai client.
type OpenAIClient interface {
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// Ensure openai.Client satisfies our interface.
var _ OpenAIClient = (*openai.Client)(nil)
