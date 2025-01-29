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
	cf "github.com/cloudflare/cloudflare-go"
	"go.uber.org/zap"
)

func validateAllowedDomains(ctx context.Context, schema *DNSRecordSchema) error {
	// No type assertion needed, you get the correct type directly
	for _, record := range schema.Records {
		// Direct access to the typed fields
		if !isAllowedDomain(record.Name) {
			return fmt.Errorf("invalid domain: %s", record.Name)
		}
	}
	return nil
}

func isAllowedDomain(domain string) bool {
	allowedDomains := getAllowedDomains()
	for _, allowedDomain := range allowedDomains {
		if strings.HasSuffix(domain, allowedDomain) {
			return true
		}
	}
	return false
}

func getAllowedDomains() []string {
	// Initialize Cloudflare API client
	api, err := cf.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		// Log error and return default domains as fallback
		logger.Log.Error("Error initializing Cloudflare client", zap.String("tag", "cloudflare-validators"), zap.Error(err))
		return nil
	}

	// Fetch zones (domains)
	zones, err := api.ListZones(context.Background())
	if err != nil {
		logger.Log.Error("Error fetching zones", zap.String("tag", "cloudflare-validators"), zap.Error(err))
		return nil
	}

	domains := make([]string, len(zones))
	for i, zone := range zones {
		domains[i] = zone.Name
	}
	return domains
}

func validateDNSRecord(ctx context.Context, schema *DNSRecordSchema) error {
	if len(schema.Records) == 0 {
		return fmt.Errorf("no records provided")
	}

	return nil
}
