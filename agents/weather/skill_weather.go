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

import "github.com/berbyte/ber-os/internal/agent"

type WeatherSkill struct {
	agent.Skill[WeatherResponseSchema]
}

var WeatherInfo = WeatherSkill{
	Skill: agent.Skill[WeatherResponseSchema]{
		Name:        "Weather Information",
		Tag:         "weather",
		Description: "Displays current weather information for a location",
		Prompt: `As a location information extractor, analyze the following text and extract:
1. The location or city name
2. Its exact latitude and longitude coordinates

Return the information in JSON format with the following fields:
- location: the name of the location
- latitude: the latitude coordinate
- longitude: the longitude coordinate`,
		Template: `### Weather in {{.location}} ({{.latitude}}, {{.longitude}}) {{if .is_day}}☀️{{else}}🌙{{end}}

🌡️ Temperature: {{.temperature}}°C
💨 Wind Speed: {{.wind_speed}} km/h
💧 Humidity: {{.humidity}}%`,
		LLMSchema: WeatherResponseSchema{},
		Hooks: agent.Hooks[WeatherResponseSchema]{
			PostLLMRequest: fetchWeatherData,
		},
	},
}
