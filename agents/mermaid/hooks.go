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

package mermaid

import (
	"context"
	"strings"
)

func postLLMRequest(ctx context.Context, response *MermaidSchema) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if !strings.HasPrefix(response.DiagramCode, "```mermaid") {
		response.DiagramCode = "```mermaid\n" + response.DiagramCode + "\n```"
	}
	return nil
}
