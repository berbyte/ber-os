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

package rest

import (
	"github.com/berbyte/ber-os/internal/adapters/github/client"

	"github.com/google/go-github/v64/github"
)

func NewIssueCommentReaction(repoOwner string, repoName string, commentID int64, reaction string) error {
	client, ctx, err := client.NewClient(repoOwner)
	if err != nil {
		return err
	}

	_, _, err = client.Reactions.CreateIssueCommentReaction(ctx, repoOwner, repoName, commentID, *github.String(reaction))
	if err != nil {
		return err
	}

	return nil
}

func NewIssueReaction(repoOwner string, repoName string, number int, reaction string) error {
	client, ctx, err := client.NewClient(repoOwner)
	if err != nil {
		return err
	}

	_, _, err = client.Reactions.CreateIssueReaction(ctx, repoOwner, repoName, number, *github.String(reaction))
	if err != nil {
		return err
	}

	return nil
}
