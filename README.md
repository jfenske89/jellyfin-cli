# jellyfin-cli

This project offers a command-line interface to Jellyfin instances.

***Notice***: This is a work in progress side project. It has limited functionality. New features will be added slowly. 

## Configuration

This tool can be configured with a config.yaml file. See the [example.config.yaml](example.config.yaml) file.

**Possible paths**:

- `/etc/jellyfin-cli/config.yaml`
- `$HOME/.config/jellyfin-cli/config.yaml`
- `./config.yaml`

## Usage

All commands follow this format:

```shell
./bin/jellyfin-cli <command> [options]
```

### Global Options

- `--output json` or `--json`: Outputs the response in JSON format instead of plain text. This can produce a lot of output, as it's returned unformatted directly from the Jellyfin API.

## Commands

### List Library Folders

Retrieve a list of all library folders grouped by collection type.

#### Command
```sh
./bin/jellyfin-cli list-library-folders
```

#### Example Output (Plain Text)
```
Library folders:
- movies:
   - Library name 1 (68094783b021abd03520a299c2c85870)
   - Library name 2 (259c2c8587096820a094783b021abd03)

- tvshows:
   - Library name 3 (a6a73b90b322532e40e44b97e68d0565)

- books:
   - Library name 4 (9515c5a412f43aeb982a99e3d6d67fec)

- music:
   - Library name 5 (7e9516eda064e319657a3c78490edccb)
```

Each section corresponds to a `CollectionType`, with library names followed by their respective `ItemId`.

### List Sessions (Users)

Retrieve a list of active user sessions, showing the username, device name, and last active time.

#### Command
```sh
./bin/jellyfin-cli list-sessions
```

#### Example Output (Plain Text)
```
Sessions:
 - User name 1 (Device name 1) 3m
 - User name 2 (Device name 2) 1m
```

The output displays:
- **User name**: The name of the logged-in user.
- **Device name**: The device being used.
- **Last active time**: Time elapsed since last activity (`m` for minutes).

### List Active Sessions

Retrieve only active sessions from the last 600 seconds (10 minutes).

#### Command
```sh
./bin/jellyfin-cli list-sessions --active
```

This filters out inactive sessions, displaying only users currently active.
