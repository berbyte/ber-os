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

package client

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"os"
	"strconv"
	"time"

	"github.com/berbyte/ber-os/internal/logger"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v64/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// Load private key from file
func loadPrivateKey() (*rsa.PrivateKey, error) {
	b64PrivateKey := os.Getenv("GH_PRIVATE_KEY")
	data, err := base64.StdEncoding.DecodeString(b64PrivateKey)
	if err != nil {
		logger.Log.Fatal("Unable to Find or Decode Private Key", zap.Error(err), zap.String("tag", "github-client-client"))
	}

	privateKeyData := []byte(data)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse private key: %w", err)
	}
	return privateKey, nil
}

// Generate JWT for GitHub App authentication
func generateJWT(appID int64, privateKey *rsa.PrivateKey) (string, error) {
	// Create the JWT claims with GitHub App specific fields
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),                       // Issued at time
		"exp": time.Now().Add(time.Minute * 10).Unix(), // Expires after 10 minutes
		"iss": appID,                                   // GitHub App ID
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with the private key
	jwtToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("Unable to sign JWT: %w", err)
	}

	return jwtToken, nil
}

// Get installation access token using the app's JWT
func getInstallationToken(client *github.Client) (string, error) {
	// List installations for the authenticated app
	installations, _, err := client.Apps.ListInstallations(context.Background(), nil)
	if err != nil {
		return "", fmt.Errorf("Failed to list installations: %w", err)
	}

	// Assume the first installation for simplicity (modify as needed)
	if len(installations) == 0 {
		return "", fmt.Errorf("No installations found for this app")
	}

	// Create an installation token for the first installation
	installationID := installations[0].GetID()
	token, _, err := client.Apps.CreateInstallationToken(context.Background(), installationID, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create installation token: %w", err)
	}

	return token.GetToken(), nil
}

func NewClient(repoOwner string) (*github.Client, context.Context, error) {
	// Define your GitHub App ID and path to the private key file
	aid, err := strconv.Atoi(os.Getenv("GH_APP_ID"))
	if err != nil {
		return nil, nil, err
	}
	appID := int64(aid)

	// Load the private key
	privateKey, err := loadPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	// Generate the JWT for authenticating as the GitHub App
	jwtToken, err := generateJWT(appID, privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Create a new GitHub client with the JWT
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: jwtToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Retrieve the installation access token
	installationToken, err := getInstallationToken(client)
	if err != nil {
		return nil, nil, err
	}

	// Create a new GitHub client with the installation access token
	ts = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: installationToken},
	)
	tc = oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	return client, ctx, nil
}
