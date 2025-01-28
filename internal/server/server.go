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

package server

import (
	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/gzip"

	"time"

	"github.com/berbyte/ber-os/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware returns a gin.HandlerFunc that logs requests using our zap logger
func LoggerMiddleware() gin.HandlerFunc {
	log := logger.GetLogger()

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Debug("gin-request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client-ip", c.ClientIP()),
		)
	}
}

func GetGinServer(appName string) *gin.Engine {
	// Set Gin to production mode to disable debug logging
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(LoggerMiddleware())
	router.Use(gin.Recovery())
	router.Use(limit.MaxAllowed(20))
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	return router
}
