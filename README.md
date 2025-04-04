# Jellyfin CLI

A command-line tool for interacting with a Jellyfin server.

***Notice***: This is a side project with limited functionality. New features will be added slowly. 

## Features

- List active sessions
- List library folders 
- View activity logs
- Search for content
- Refresh library

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/jfenske89/jellyfin-cli.git
cd jellyfin-cli
```

2. Build the application:
```bash
go build -o ./bin
```

3. (Optional) Install to your system:
```bash
go install
```

## Configuration

Create a configuration file at one of the following locations:
- `./config.yaml` (current directory)
- `$HOME/.config/jellyfin-cli/config.yaml`
- `/etc/jellyfin-cli/config.yaml`

Example configuration:

```yaml
logging:
  # level is the logging level: DEBUG, INFO, WARN, ERROR
  level: INFO

api:
  # base_url is the base URL to a Jellyfin instance 
  base_url: http://127.0.0.1:8096
  # token is an API token for authentication
  token: your-api-token-here
  # insecure can permit insecure SSL requests
  insecure: false
```

## Getting an API Token

To use this CLI, you need to generate an API token from your Jellyfin server:

1. Go to your Jellyfin dashboard
2. Navigate to Admin > Dashboard > Advanced > API Keys
3. Create a new API key
4. Copy the key to your config file

## Usage

### General Help

```bash
jellyfin-cli help
```

### List Sessions

List all sessions:
```bash
jellyfin-cli sessions
```

List only active sessions:
```bash
jellyfin-cli sessions --active
```

### List Libraries

List all library folders:
```bash
jellyfin-cli libraries
```

### Refresh Library

Trigger a library refresh:
```bash
jellyfin-cli libraries refresh
```

### Activity Logs

View recent activity logs:
```bash
jellyfin-cli activity
```

Limit the number of logs:
```bash
jellyfin-cli activity --limit 20
```

### Search

Search for content:
```bash
jellyfin-cli search "star wars"
```

Filter by content type:
```bash
jellyfin-cli search "star wars" --type Movie
```

Limit search results:
```bash
jellyfin-cli search "star wars" --limit 5
```

### JSON Output

Any command can output JSON by adding the `--json` flag:

```bash
jellyfin-cli sessions --json
jellyfin-cli libraries --json
jellyfin-cli activity --json
jellyfin-cli search "star wars" --json
```
