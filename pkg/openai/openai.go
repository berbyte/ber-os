// Copyright 2025 BER - ber.run
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	oai "github.com/openai/openai-go"
)

// Client wraps the OpenAI client with additional functionality
type Client struct {
	client *oai.Client
}

// NewClient creates a new OpenAI client wrapper
func NewClient() *Client {
	return &Client{
		client: oai.NewClient(),
	}
}

// GenerateSchema creates a JSON schema for any struct type
func GenerateSchema(T interface{}) interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		// Add support for embedded structs
		ExpandedStruct: true,
		// Add support for recursive types
		DoNotReference: true,
	}
	schema := reflector.Reflect(T)
	return schema
}

// Role type for chat messages
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// ChatMessage represents a message in the chat thread
type ChatMessage struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// QueryWithSchema sends a query to OpenAI and returns structured data based on the provided schema type
func QueryWithSchema(ctx context.Context, client *Client, messages []ChatMessage, resp_schema interface{}) (interface{}, error) {
	schema := GenerateSchema(resp_schema)

	schemaParam := oai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        oai.F("structured_response"),
		Description: oai.F("Structured response based on schema"),
		Schema:      oai.F(schema),
		Strict:      oai.Bool(true),
	}

	// Convert ChatMessage slice to OpenAI message format
	apiMessages := make([]oai.ChatCompletionMessageParamUnion, len(messages))
	for i, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			apiMessages[i] = oai.SystemMessage(msg.Content)
		case RoleUser:
			apiMessages[i] = oai.UserMessage(msg.Content)
		case RoleAssistant:
			apiMessages[i] = oai.AssistantMessage(msg.Content)
		}
	}

	chat, err := client.client.Chat.Completions.New(ctx, oai.ChatCompletionNewParams{
		Messages: oai.F(apiMessages),
		ResponseFormat: oai.F[oai.ChatCompletionNewParamsResponseFormatUnion](
			oai.ResponseFormatJSONSchemaParam{
				Type:       oai.F(oai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: oai.F(schemaParam),
			},
		),
		Model: oai.F(oai.ChatModelGPT4oMini),
	})
	if err != nil {
		return nil, fmt.Errorf("[PKG-OAI] OpenAI API error: %w", err)
	}

	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &resp_schema)
	if err != nil {
		return nil, fmt.Errorf("[PKG-OAI] Failed to parse response: %w", err)
	}

	return resp_schema, nil
}
