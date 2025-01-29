# Cloudflare Agent

A BER agent for managing DNS records through Cloudflare's API.

## Features

- DNS record creation (A, CNAME, MX, TXT)
- Automated record validation
- Domain verification
- Cloudflare proxy support

## Prerequisites

- Cloudflare API Token with DNS management permissions
- Environment variable: `CLOUDFLARE_API_TOKEN`


## Usage

```
export CLOUDFLARE_API_TOKEN=""

go run . webhook --debug
```


The agent provides DNS management capabilities through the following skill:
```go
DNSManagement:
  - Record creation
  - Syntax validation
  - Domain verification
  - TTL configuration
  - Proxy status management
```

## Validation

The agent performs several validations:
- DNS record syntax
- Allowed domains (based on Cloudflare account zones)
- Content format per record type
