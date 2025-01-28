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
	"fmt"
	"math"

	oai "github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

// GetEmbedding retrieves the vector embedding for a given text input using OpenAI's text-embedding-ada-002 model
func (c *Client) GetEmbedding(ctx context.Context, input string) ([]float32, error) {
	// Create embedding request parameters
	embedParams := oai.EmbeddingNewParams{
		Model: oai.F(oai.EmbeddingModelTextEmbeddingAda002),
		Input: oai.F[oai.EmbeddingNewParamsInputUnion](shared.UnionString(input)),
	}

	// Call OpenAI embedding API
	resp, err := c.client.Embeddings.New(ctx, embedParams)
	if err != nil {
		return nil, fmt.Errorf("[PKG-EMD] Failed to get embedding: %w", err)
	}

	// Check if we got any embeddings back
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("[PKG-EMD] No embeddings returned from API")
	}

	// Convert []float64 to []float32
	embedding := make([]float32, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		embedding[i] = float32(v)
	}
	return embedding, nil
}

// CosineSimilarity calculates the cosine similarity between two embedding vectors
func CosineSimilarity(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("[PKG-EMD] Vectors must have same length, got %d and %d", len(a), len(b))
	}

	var dotProduct float32
	var normA float32
	var normB float32

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	// Avoid division by zero
	if normA == 0 || normB == 0 {
		return 0, fmt.Errorf("[PKG-EMD] Vector norm cannot be zero")
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB)))), nil
}
