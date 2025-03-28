# ntgrrc

ntgrrc (Netgear Remote Control) a command line (CLI) tool to manage Netgear managed plus switches 300 series.

Since Netgear does not offer a REST API, this tool uses web scraping techniques to
manage configuration and show status of the switch.

This tool is still very limited in its features and more testers and **contributors
are very welcome**.

### Build Status

[![go test](https://github.com/nitram509/ntgrrc/actions/workflows/go-test.yml/badge.svg)](https://github.com/nitram509/ntgrrc/actions/workflows/go-test.yml)
[![codecov](https://codecov.io/gh/nitram509/ntgrrc/branch/main/graph/badge.svg?token=8LVPP8JVKY)](https://codecov.io/gh/nitram509/ntgrrc)

### Supported firmware versions

A green tick means it was successfully tested. 

| Firmware  | GS305EP(P) | GS308EP  | GS308EPP | GS316EP(P) |
|-----------|------------|----------|----------|------------|
| V1.0.0.8  | ✅          | (likely) | ✅        | -          |
| v1.0.0.10 | ✅          | (likely) | ✅        | -          |
| V1.0.1.1  | (likely)   | ✅        | (likely) | -          |
| V1.0.3.4  | -          | -        | -        | ✅          |
| V1.0.3.7  | -          | -        | -        | ✅          |
| V1.0.4.4  | -          | -        | -        | ✅          |

Legend: \
✅ = successfully tested \
`-`  = not available \

## download & installation

This tool is build with the Go programming language
and pre-build binaries for Windows, Linux, and MacOSX are available for [download](https://github.com/nitram509/ntgrrc/releases).

Just download the fitting binary for your operating system and put it somewhere in your PATH.

## usage

### help

```shell
ntgrrc --help-all
```

<!-- MARKDOWN-AUTO-DOCS:START (CODE:src=./help.txt) -->
<!-- The below code snippet is automatically added from ./help.txt -->
```txt
Usage: ntgrrc <command> [flags]

Flags:
  -h, --help                  Show context-sensitive help.
      --help-all              advanced/full help
  -v, --verbose               verbose log messages
  -q, --quiet                 no log messages
  -f, --output-format="md"    what output format to use [md, json]
  -d, --token-dir=""          directory to store login tokens

Commands:
  version [flags]
    show version

  login --address=STRING [flags]
    create a session for further commands (requires admin console password)

  poe status --address=STRING
    show current PoE status for all ports

  poe settings --address=STRING
    show current PoE settings for all ports

  poe set --address=STRING --port=PORT,... [flags]
    set new PoE settings per each PORT number

  poe cycle --address=STRING --port=PORT,...
    power cycle one or more PoE ports

  port settings --address=STRING
    show switch port settings

  port set --address=STRING --port=PORT,... [flags]
    set properties for a port number

  debug-report --address=STRING [flags]
    show information from the switch communication, useful for supporting
    development and bug fixes

Run "ntgrrc <command> --help" for more information on a command.
```
<!-- MARKDOWN-AUTO-DOCS:END -->

### login

For better performance, **login first**.
The login action will store a token to a file called ```$TEMP/.config/ntgrrc/token-12345678```
and thus subsequent actions will use it and are authenticated.

Note: if you have multiple Netgear switches, ntgrrc **supports multiple parallel tokens**/sessions,
because the token file's name is derived from the provided ```--address``` device name.

```shell
ntgrrc login --address gs305ep --password secret
```

### show port settings

Once a session is created, you can fetch port settings.

#### Settings 

The switch's port settings are printed in Markdown table format.
This means, separated by | (pipe) and optional suffixes with blanks.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc port settings --address gs305ep```

```markdown
| Port ID | Port Name | Speed | Ingress Limit | Egress Limit | Flow Control | Port Status | Link Speed |
|---------|-----------|-------|---------------|--------------|--------------|-------------|------------|
| 1       | XYZ       | Auto  | No Limit      | No Limit     | Off          | AVAILABLE   | No Speed   |
| 2       |           | Auto  | No Limit      | No Limit     | On           | CONNECTED   | 100M Half  |
| 3       |           | Auto  | No Limit      | No Limit     | On           | DISABLED    | No Speed   |
| 4       |           | Auto  | 1 Mbit/s      | No Limit     | On           | AVAILABLE   | No Speed   |
```

### set port settings

ntgrrc is able to set various parameters on switch port(s).

#### Port Name

To change the port name (within the switch's limit of 1-16 characters), pass the name using `-n` and the desired name in quotes. More than one port number can be provided.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc port set -p 1 -n 'port #1' --address gs305ep```

```markdown
| Port ID | Port Name | Speed | Ingress Limit | Egress Limit | Flow Control |
|---------|-----------|-------|---------------|--------------|--------------|
| 1       | port #1   | Auto  | No Limit      | No Limit     | Off          |

```

To clear the set name, supply an empty, but quoted string. More than one port number can be provided.

```ntgrrc port set -p 1 -n '' --address gs305ep```

```markdown
| Port ID | Port Name | Speed | Ingress Limit | Egress Limit | Flow Control |
|---------|-----------|-------|---------------|--------------|--------------|
| 1       |           | Auto  | No Limit      | No Limit     | Off          |
```

#### Speed

To change the port speed, use `-s` and the desired speed ('100M full', '100M half', '10M full', '10M half', 'Auto', 'Disable') in quotes. More than one port number can be provided.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc port set -p 1 -s '100M half' --address gs305ep```

```markdown
| Port ID | Port Name | Speed     | Ingress Limit | Egress Limit | Flow Control |
|---------|-----------|-----------|---------------|--------------|--------------|
| 1       |           | 100M half | No Limit      | No Limit     | Off          |
```

#### In Rate Limit

To change the in rate limit, use `-i` and the desired rate limit ('1 Mbit/s', '128 Mbit/s', '16 Mbit/s', '2 Mbit/s', '256 Mbit/s', '32 Mbit/s', '4 Mbit/s', '512 Kbit/s', '512 Mbit/s', '64 Mbit/s', '8 Mbit/s', 'No Limit') in quotes. More than one port number can be provided.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc port set -p 1 -i '16 Mbit/s' --address gs305ep```

```markdown
| Port ID | Port Name | Speed | Ingress Limit | Egress Limit | Flow Control |
|---------|-----------|-------|---------------|--------------|--------------|
| 1       |           | Auto  | 16 Mbit/s     | No Limit     | Off          |
```

#### Out Rate Limit

To change the out rate limit, use `-o` and the desired rate limit ('1 Mbit/s', '128 Mbit/s', '16 Mbit/s', '2 Mbit/s', '256 Mbit/s', '32 Mbit/s', '4 Mbit/s', '512 Kbit/s', '512 Mbit/s', '64 Mbit/s', '8 Mbit/s', 'No Limit') in quotes. More than one port number can be provided.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc port set -p 1 -o '16 Mbit/s' --address gs305ep```

```markdown
| Port ID | Port Name | Speed | Ingress Limit | Egress Limit | Flow Control |
|---------|-----------|-------|---------------|--------------|--------------|
| 1       |           | Auto  | 16 Mbit/s     | 16 Mbit/s    | Off          |
```

#### Flow Control

To change the flow control setting for a port, use `--flow-control` and the desired setting ('On', 'Off') in quotes. More than one port number can be provided.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc port set -p 1 --flow-control 'On' --address gs305epp```

```markdown
ntgrrc port set -p 1 --flow-control 'On' --address test
| Port ID | Port Name | Speed | Ingress Limit | Egress Limit | Flow Control |
|---------|-----------|-------|---------------|--------------|--------------|
| 1       |           | Auto  | 16 Mbit/s     | 16 Mbit/s    | On           |
```

### show Power Over Ethernet (POE)

Once a session is created, you can fetch POE settings and status.

#### Settings

The switch's PoE settings are printed in Markdown table format.
This means, separated by | (pipe) and optional suffixes with blanks.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc poe settings --address gs316ep```

```markdown
| Port ID | Port Name        | Port Power | Mode        | Priority | Limit Type | Limit (W) | Type                 | Longer Detection Time |
|---------|------------------|------------|-------------|----------|------------|-----------|----------------------|-----------------------|
| 1       | AGER 31 SUR Tech | enabled    | Legacy      | High     | User       | 30.0      | IEEE802              | Disable               |
| 2       | foobar           | enabled    | 802.3at     | Low      | User       | 30.0      | IEEE802              | Disable               |
| 3       | zzz              | enabled    | 802.3at     | Low      | User       | 30.0      | 4pt 802.3af + Legacy | Disable               |
| 4       | uuu              | enabled    | 802.3at     | Low      | User       | 30.0      | IEEE802              | Disable               |
```

#### Status

The switch's POE status are printed in Markdown table format.
This means, separated by | (pipe) and optional suffixes with blanks.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc poe status --address gs305ep```

```markdown
| Port ID | Port Name | Status           | PortPwr class | Voltage (V) | Current (mA) | PortPwr (W) | Temp. (°C) | Error status |
|---------|-----------|------------------|---------------|-------------|--------------|-------------|------------|--------------|
| 1       | Camera    | Delivering Power | 0             | 53          | 82           | 4.40        | 30         | No Error     |
| 2       |           | Searching        |               | 0           | 0            | 0.00        | 30         | No Error     |
| 3       |           | Searching        |               | 0           | 0            | 0.00        | 30         | No Error     |
| 4       |           | Searching        |               | 0           | 0            | 0.00        | 30         | No Error     |
```

### set Power Over Ethernet (POE)

ntgrrc is able to set various parameters on PoE port(s).

#### Port Power

To enable or disable port power, pass the port number using `-p` and `--power enable` to enable power or `--power disable` to disable power. More than one port number can be provided.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc poe set -p 3 -p 4 --power enable --address gs305ep```

```markdown
| Port ID | Port Name | Port Power | Mode        | Priority | Limit Type | Limit (W) | Type    | Longer Detection Time |
|---------|-----------|------------|-------------|----------|------------|-----------|---------|-----------------------|
| 3       |           | enabled    | Legacy      | High     | User       | 30.0      | IEEE802 | Disable               |
| 4       |           | enabled    | 802.3at     | Low      | User       | 30.0      | IEEE802 | Disable               |

```

#### Port Power Mode

To change the port power mode, pass the port number using `-p` and `--mode` with the desired power mode (802.3af, legacy, pre-802.3at, 802.3at). More than one port number can be provided.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc poe set -p 3 -p 5 --mode legacy --address gs305ep```

```markdown
| Port ID | Port Name | Port Power | Mode        | Priority | Limit Type | Limit (W) | Type    | Longer Detection Time |
|---------|-----------|------------|-------------|----------|------------|-----------|---------|-----------------------|
| 3       |           | enabled    | Legacy      | High     | User       | 30.0      | IEEE802 | Disable               |
| 4       |           | enabled    | Legacy      | Low      | User       | 30.0      | IEEE802 | Disable               |
```

#### Port Priority

To change port priority, pass the port number using `-p` and `--priority` with the desired priority (low, high, critical). More than one port number can be provided.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc poe set -p 3 -p 5 --priority critical --address gs305ep```

```markdown
| Port ID | Port Name | Port Power | Mode        | Priority | Limit Type | Limit (W) | Type    | Longer Detection Time |
|---------|-----------|------------|-------------|----------|------------|-----------|---------|-----------------------|
| 3       |           | enabled    | Legacy      | critical | User       | 30.0      | IEEE802 | Disable               |
| 5       |           | enabled    | Legacy      | critical | User       | 30.0      | IEEE802 | Disable               |
```

#### Power Limit

To change the power limit for a port, pass the port number using `-p` and `--pwr-limit` with the desired limit. More than one port number can be provided. 

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc poe set -p 3 -p 5 --pwr-limit 5 --address gs305ep```

```markdown
| Port ID | Port Name | Port Power | Mode        | Priority | Limit Type | Limit (W) | Type    | Longer Detection Time |
|---------|-----------|------------|-------------|----------|------------|-----------|---------|-----------------------|
| 3       |           | enabled    | Legacy      | critical | User       | 5.0       | IEEE802 | Disable               |
| 5       |           | enabled    | Legacy      | critical | User       | 5.0       | IEEE802 | Disable               |
```

#### Power Limit Type

To change the power limit type for a port, pass the port number using `-p` and `--limit-type` with the desired limit type (none, class, user). More than one port number can be provided.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc poe set -p 3 -p 5 --limit-type class --address gs305ep```

```markdown
| Port ID | Port Name | Port Power | Mode        | Priority | Limit Type | Limit (W) | Type    | Longer Detection Time |
|---------|-----------|------------|-------------|----------|------------|-----------|---------|-----------------------|
| 3       |           | enabled    | Legacy      | critical | class      | 30.0      | IEEE802 | Disable               |
| 5       |           | enabled    | Legacy      | critical | class      | 30.0      | IEEE802 | Disable               |
```

#### Detection type

To change the detection type for a port, pass the port number using `-p` and `--detect-type` with the desired detection type. More than one port can be provided.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc poe set -p 3 -p 5 --detect-type "4pt 802.3af + Legacy" -a gs305ep```

```markdown
| Port ID | Port Name | Port Power | Mode        | Priority | Limit Type | Limit (W) | Type                 | Longer Detection Time |
|---------|-----------|------------|-------------|----------|------------|-----------|----------------------|-----------------------|
| 3       |           | enabled    | Legacy      | critical | User       | 30.0      | 4pt 802.3af + Legacy | Disable               |
| 5       |           | enabled    | Legacy      | critical | User       | 30.0      | 4pt 802.3af + Legacy | Disable               |
```

#### cycle Power Over Ethernet (POE)

ntgrrc is able to power cycle one or more PoE ports.

Use the ```--output-format=json``` flag, to get JSON output instead.

```ntgrrc poe cycle -p 3 -p 5 --address gs305ep```
Alternativly, you can achieve the same request with this
```ntgrrc poe cycle --address gs305ep --port=3,5```

```markdown
| Port ID | Port Name        | Status           | PortPwr class | Voltage (V) | Current (mA) | PortPwr (W) | Temp. (°C) | Error status |
|---------|------------------|------------------|---------------|-------------|--------------|-------------|------------|--------------|
| 3       | Camera           | Delivering Power |               | 54          | 24           | 1.30        | 30         | No Error     |
| 5       | Sensor           | Searching        |               | 0           | 0            | 0.00        | 30         | Power Denied |
```
