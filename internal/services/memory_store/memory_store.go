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

package memory_store

import (
	"strings"

	"github.com/berbyte/ber-os/internal/logger"
	"go.uber.org/zap"
)

// MemoryStore provides a simple key-value storage mechanism
type MemoryStore struct {
	data map[string]interface{}
}

var globalStore *MemoryStore

// NewMemoryStore creates a new instance of MemoryStore
func NewMemoryStore() *MemoryStore {
	if globalStore != nil {
		logger.GetLogger().Info("Returning existing global store",
			zap.Int("item_count", len(globalStore.data)))
		return globalStore
	}
	logger.GetLogger().Info("Creating new MemoryStore")
	globalStore = &MemoryStore{
		data: make(map[string]interface{}),
	}
	return globalStore
}

// Get retrieves a value from the store by key
func (m *MemoryStore) Get(key string) (interface{}, bool) {
	value, exists := m.data[key]
	logger.GetLogger().Debug("Getting key from store",
		zap.String("key", key),
		zap.Bool("exists", exists))
	return value, exists
}

// Set stores a value in the store with the given key
func (m *MemoryStore) Set(key string, value interface{}) {
	logger.GetLogger().Debug("Setting key in store",
		zap.String("key", key),
		zap.Any("value", value))
	m.data[key] = value
}

// Delete removes a key-value pair from the store
func (m *MemoryStore) Delete(key string) {
	logger.GetLogger().Debug("Deleting key from store",
		zap.String("key", key))
	delete(m.data, key)
}

// Len returns the number of items in the memory store
func (m *MemoryStore) Len() int {
	length := len(m.data)
	logger.GetLogger().Debug("Getting store length",
		zap.Int("length", length))
	return length
}

// GetAllData returns all data in the memory store
func (m *MemoryStore) GetAllData() map[string]interface{} {
	logger.GetLogger().Debug("Getting all data from store",
		zap.Int("item_count", len(m.data)))
	return m.data
}

// GetKeysWithSubstring returns all keys that contain the given substring
func (m *MemoryStore) GetKeysWithSubstring(substring string) []string {
	var matches []string
	for key := range m.data {
		if strings.Contains(key, substring) {
			matches = append(matches, key)
		}
	}
	logger.GetLogger().Debug("Found keys containing substring",
		zap.String("substring", substring),
		zap.Strings("matches", matches),
		zap.Int("match_count", len(matches)))
	return matches
}
