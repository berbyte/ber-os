<a id="readme-top"></a>

> [!note]
>
> 🚧 **Active Development Phase** 🚧
>
> BER is under active development. Features and APIs are subject to change as we work toward a stable release. We welcome early adopters and value your feedback.
>
> [Give Feedback](https://github.com/berbyte/ber-os/discussions/new?category=feedback)


<div align="center">
    <img src="https://rtfm.ber.run/ber-intro.png" alt="BER: LLM SuperGlue">
  <h1 align="center">BER: LLM SuperGlue v0.1</h1>
  <p align="center">
    <a href="https://rtfm.ber.run"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="#demo">View Demo</a>
    &middot;
    <a href="https://github.com/berbyte/ber-os/issues/new">Report a Bug</a>
    &middot;
    <a href="#getting-started">Quickstart Guide</a>
  </p>

</div>

## What is this all about?

BER is the glue enabling your own AI-powered workflows. With BER, you can employ the power of AI in your problem domains. Connect your workflows to any interface and collaborate with any API—all in natural language.

## Demo

<div align="center">
  <img src="https://rtfm.ber.run/demo.gif" alt="BER Demo">
</div>


> [!tip]
> **Want to see BER in action?**
>
> Check out our [use cases documentation](https://rtfm.ber.run/getting-started/usecases/) to explore examples of how BER can help you.
>

## Why BER?

- 🛑 Manual workflow automation is slow and error-prone, wasting valuable time.
- ✅ BER helps engineers, developers, and businesses to integrate tools and APIs using the power of LLMs.

### Key Features

- **Superglue for Your Tools**: Connect ITSM, chat, ticketing, git platforms, and more with any API, database, or cloud provider.
- **Universal Assistance**: Customize BER in new roles, access them from your closest lying interface.
- **Highly Customizable Workflows**: Build workflows tailored to your organization's unique needs.
- **Control over AI**: Use the combination of Skills, Hooks, Actions, and Validators to control conversational LLM-type AI models.
- **Developer-First Design**: Built with engineers in mind, offering intuitive APIs and straightforward documentation.
- **Open-Source Freedom**: Fully open-source under an Apache license—no vendor lock-in, complete transparency, and a community-driven ecosystem.


## High-Level Diagram

<div align="center">
    <a href="https://rtfm.ber.run">
    <img src="https://rtfm.ber.run/diagrams/ber-intro-splash.svg" alt="BER: LLM SuperGlue">
  </a>
</div>

> [!note]
> If you like drawings, check out our [documentation website](https://rtfm.ber.run) for [architecture diagrams](https://rtfm.ber.run/getting-started/), [workflows](https://rtfm.ber.run/concepts/agent/), and [more](https://rtfm.ber.run/concepts/adapter/)!


<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- Agent -->
## Available Agents
BER comes with several built-in agents to help you get started:

- **[Cloudflare Agent](agents/cloudflare)**: Helps manage DNS records
  - Skills: DNS Management
- **[Mermaid Diagram Agent](agents/mermaid)**: Helps create various types of Mermaid diagrams
  - Skills: Flowchart Diagrams, Pie Charts, Sequence Diagrams
- **[Product Owner Agent](agents/product)**: Helps with story refinement, acceptance criteria, and estimation
  - Skills: Acceptance Criteria, Story Estimation, Story Refinement
- **[Weather Visualization Agent](agents/weather)**: Helps visualize and analyze weather data using charts and diagrams
  - Skills: Weather Trend Charts

Each agent can be customized and extended to meet your specific needs. For detailed information about configuring and using these agents, visit our [documentation](https://rtfm.ber.run/concepts/agent/).

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## Getting Started
1. Clone the repository:

```
git clone git@github.com:berbyte/ber-os.git
cd ber-os
```

2. Install the dependencies:

```
go mod tidy
```

### TUI Usage
> [!CAUTION]
> The TUI adapter is currently in an experimental state. We recommend using the [GitHub adapter](#github-application-usage) which provides a more stable and feature-complete experience.

1. Set the environment variable

```
export OPENAI_API_KEY=""
```

2. Run the TUI:

```
go run . tui
```

### GitHub Application Usage
0. Create a GitHub App by following our [documentation guide](https://rtfm.ber.run/guides/howto-adapter-github-install/). This will provide you with the required credentials for the next steps.

1. Set the environment variables
```
export GH_APP_ID=""
export GH_PRIVATE_KEY="" # base64 decoded pem
export GH_WEBHOOK_SECRET=""

export OPENAI_API_KEY=""
```

2. Start ngrok:

```
ngrok http http://localhost:8080
```

3. Run the application:
```
go run . webhook --debug
```

For detailed GitHub adapter usage instructions, please visit our [GitHub Adapter Tutorial](https://rtfm.ber.run/tutorials/github/).


### Build Your First Agent
Now that you have the environment set up, you're ready to build your first BERAgent! Check out our [Agent Building Tutorial](https://rtfm.ber.run/tutorials/agent/) to get started with creating custom agents for your specific use cases.



<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- CONTRIBUTING -->
## Contributing
Any contributions you make are **greatly appreciated**. We would love to hear your feedback - feel free to [open a new discussion](https://github.com/berbyte/ber-os/discussions/new?category=feedback)!

Please read our [Contributing Guidelines](.github/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

<!-- SECURITY -->
## Security
We take security seriously. If you believe you have found a security vulnerability, please report it to us as described in our [Security Policy](.github/SECURITY.md).

<!-- CODE OF CONDUCT -->
## Code of Conduct
This project and everyone participating in it is governed by our [Code of Conduct](.github/CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

<!-- LICENSE -->
## License
Distributed under the Apache-2.0 license. See [`LICENSE.txt`](LICENSE.txt) for more information.

<!-- CONTACT -->
## Contact
- BER - github@ber.run
- Project Link: [https://github.com/berbyte/ber-os](https://github.com/berbyte/ber-os)
- Documentation: [https://rtfm.ber.run](https://rtfm.ber.run)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
