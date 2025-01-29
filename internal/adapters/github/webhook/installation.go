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

package webhook

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/berbyte/ber-os/internal/adapters/github/rest"
	"github.com/berbyte/ber-os/internal/logger"
	g "github.com/google/go-github/v64/github"
	"go.uber.org/zap"
)

func consumeInstallEvent(payload []byte) error {
	var p g.InstallationEvent
	err := json.Unmarshal(payload, &p)
	if err != nil {
		logger.Log.Error("Failed to Unmarshal Installation Event",
			zap.String("tag", "github-webhook-installation"), zap.Error(err))
		return err
	}

	if p.Action != nil && *p.Action == "created" { // Create event
		if len(p.Repositories) > 0 && p.Repositories[0] != nil {
			repo := strings.Split(*p.Repositories[0].FullName, "/")

			welcomeMessage := `

# Welcome to BER 🎉!

We're excited to have you on board!

### Available BERAgents:


### 📊 **Mermaid Agent** - Visualize your architecture and flows

You can try these prompts:


` +
				"```\n@ber Create a sequence diagram for user authentication flow\n```\n\n" +
				"```\n@ber Generate a flowchart for the CI/CD pipeline\n```\n\n" +
				`

### 🌐 **Cloudflare Agent** - Create DNS entries

Try these prompts:


` +

				"```\n@ber create new record for demo.ber.run -> google.com\n```\n\n" +
				"```\n@ber changed my mind, point it to google.com and lower the ttl\n```\n\n" +
				`

Feel free to experiment, simply create a new comment on this issue and mention BER using @ber.

📚 For detailed documentation and more examples, visit: https://rtfm.ber.run

Happy hacking! 🚀
dOMiNiS

> [!NOTE]
>
> 💭 Have feedback? We'd love to hear from you! Send your thoughts to github@ber.run
`

			err := rest.NewIssue(
				repo[0],
				repo[1],
				"Welcome to BER!",
				welcomeMessage,
				*p.Sender.Login,
			)
			if err != nil {
				logger.Log.Error("Failed to Post Reaction on Comment", zap.String("tag", "github-webhook-issues"), zap.Error(err))
				return nil
			}
			logger.Log.Info("BERAdapter for Github Is Installed",
				zap.String("tag", "github-webhook-installation"),
				zap.String("repository", fmt.Sprintf("%s/%s", repo[0], repo[1])))
		}
	}
	return nil
}
