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

package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type weatherResponse struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
		WindSpeed   float64 `json:"wind_speed_10m"`
		RelHumidity int     `json:"relative_humidity_2m"`
		IsDay       int     `json:"is_day"`
	} `json:"current"`
}

func fetchWeatherData(ctx context.Context, response *WeatherResponseSchema) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Get weather data using the coordinates
	weatherURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.6f&longitude=%.6f&current=temperature_2m,relative_humidity_2m,wind_speed_10m,is_day",
		response.Latitude, response.Longitude)

	req, err := http.NewRequestWithContext(ctx, "GET", weatherURL, nil)
	if err != nil {
		return fmt.Errorf("creating weather request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetching weather data: %w", err)
	}
	defer resp.Body.Close()

	var weatherResp weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return fmt.Errorf("decoding weather response: %w", err)
	}

	// Update the response with real weather data
	response.Temperature = weatherResp.Current.Temperature
	response.WindSpeed = weatherResp.Current.WindSpeed
	response.Humidity = weatherResp.Current.RelHumidity
	response.IsDay = weatherResp.Current.IsDay == 1

	return nil
}
