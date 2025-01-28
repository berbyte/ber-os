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

func NewIssue(repoOwner string, repoName string, title string, body string, assignee string) error {
	client, ctx, err := client.NewClient(repoOwner)
	if err != nil {
		return err
	}

	// Create a new issue comment
	issue := &github.IssueRequest{
		Title:    github.String(title),
		Body:     github.String(body),
		Assignee: github.String(assignee),
	}

	// Post the comment to the specified issue
	_, _, err = client.Issues.Create(ctx, repoOwner, repoName, issue)
	if err != nil {
		return err
	}

	return nil
}

func GetIssue(repoOwner string, repoName string, issueNumber int) (*github.Issue, error) {
	client, ctx, err := client.NewClient(repoOwner)
	if err != nil {
		return nil, err
	}

	issue, _, err := client.Issues.Get(ctx, repoOwner, repoName, issueNumber)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func GetIssueLabels(repoOwner string, repoName string, issueNumber int) ([]*github.Label, error) {
	client, ctx, err := client.NewClient(repoOwner)
	if err != nil {
		return nil, err
	}

	issue, _, err := client.Issues.Get(ctx, repoOwner, repoName, issueNumber)
	if err != nil {
		return nil, err
	}

	return issue.Labels, nil
}
