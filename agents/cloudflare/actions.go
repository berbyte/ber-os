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

package cloudflare

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/berbyte/ber-os/internal/logger"
	"github.com/cloudflare/cloudflare-go"
	"go.uber.org/zap"
)

func createRecords(ctx context.Context, data *DNSRecordSchema) error {
	// Initialize Cloudflare API client
	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudflare client: %w", err)
	}

	// Extract domain from record name by taking last two segments
	parts := strings.Split(data.Records[0].Name, ".")
	domain := strings.Join(parts[len(parts)-2:], ".")
	logger.Log.Info("Extracted domain from record name",
		zap.String("domain", domain),
		zap.String("record name", data.Records[0].Name),
		zap.String("tag", "cloudflare-actions"))

	// Get zone ID for domain
	zoneID, err := getZoneID(ctx, api, domain)
	if err != nil {
		return fmt.Errorf("failed to get zone ID: %w", err)
	}
	rc := cloudflare.ZoneIdentifier(zoneID)

	// Create DNS record
	record := cloudflare.CreateDNSRecordParams{
		Type:    data.Records[0].Type,
		Name:    data.Records[0].Name,
		Content: data.Records[0].Content[0],
		TTL:     data.Records[0].TTL,
		Proxied: &data.Records[0].Proxied,
	}

	_, err = api.CreateDNSRecord(ctx, rc, record)
	if err != nil {
		return fmt.Errorf("failed to create DNS record: %w", err)
	}

	logger.Log.Info("Successfully created DNS record",
		zap.String("record name", data.Records[0].Name),
		zap.String("tag", "cloudflare-actions"))
	return nil
}

// getZoneID retrieves the zone ID for a given domain
func getZoneID(ctx context.Context, api *cloudflare.API, domain string) (string, error) {
	zones, err := api.ListZones(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list zones: %w", err)
	}

	for _, zone := range zones {
		if zone.Name == domain {
			return zone.ID, nil
		}
	}
	return "", fmt.Errorf("no zone found for domain: %s", domain)
}
