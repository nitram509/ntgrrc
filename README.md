# ntgrrc

ntgrrc (Netgear Remote Control) a command line (CLI) tool to manage Netgear managed plus switches 300 series.

Since Netgear does not offer a REST API, this tool uses web scraping techniques to
manage configuration and show status of the switch.

This tool is still very limited in its features and more testers and contributors
are very welcome.

### Supported firmware versions

| Firmware  | GS305EP | GS305EPP |
|-----------|---------|----------|
| V1.0.0.8  | ✅       | ✅        |
| v1.0.0.10 | ✅       | ✅        |


## download & installation

This tool is build with the Go programming language
and pre-build binaries for Windows, Linux, and MacOSX are available for [download](https://github.com/nitram509/ntgrrc/releases).

Just download the fitting binary for your operating system und put it somewhere in your PATH.

## usage

### login

For better performance, **login first**.
The login action will store a token to a file called ```~/.config/ntgrrc/token-12345678```
and thus subsequent actions will use it and are authenticated.

Note: if you have multiple Netgear switches, ntgrrc **supports multiple parallel tokens**/sessions,
because the token file's name is derived from the provided ```--address``` device name.

```shell
ntgrrc login --address gs305ep --password secret
```


### show Power Over Ethernet (POE)

Once a session is created, you can fetch POE settings and status.

#### Settings 

```ntgrrc poe settings --address gs305ep```

```text
Port ID | Port Power |        Mode | Priority | Limit Type | Limit (W) |                 Type
      1 |   disabled |     802.3at |      low |       user |      30.0 |             IEEE 802
      2 |    enabled |     802.3at |      low |       user |      30.0 |             IEEE 802
      3 |    enabled |     802.3at |      low |       user |      30.0 |             IEEE 802
      4 |    enabled |     802.3at |      low |       user |      30.0 |             IEEE 802
```
#### Status

```ntgrrc poe status --address gs305ep```

```text
Port ID |       Status | Power class | Voltage (V) | Current (mA) | Power (W) | Temperature (°C) | Error status
      1 |     Disabled |             |           0 |            0 |  0.000000 |               32 | No Error
      2 |    Searching |             |           0 |            0 |  0.000000 |               33 | No Error
      3 |    Searching |             |           0 |            0 |  0.000000 |               33 | No Error
      4 |    Searching |             |           0 |            0 |  0.000000 |               33 | No Error
```
