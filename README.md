# ntgrrc

ntgrrc (Netgear Remote Control) a command line (CLI) tool to manage Netgear GS3xx switch series.

Since Netgear does not offer a REST API, this tool uses web scraping techniques to
manage configuration and show status of the switch.

This tool is still very limited in its features and more testers and contributors
are very welcome.

Tested with GS305EP, v1.0.8 firmware.

## download-installation

This tool is build with the Go programming language
and pre-build binaries for Windows, Linux, and MacOSX are available for [download](https://github.com/nitram509/ntgrrc/releases).


## usage

### login

For better performance, login a subsequent actions are separated.
So, please create a session via login first.

```shell
ntgrrc login --address gs305ep --password secret
```


### show Power Over Ethernet (POE) status

Once a session is created, you can fetch POE status.

```shell
ntgrrc poe status --address gs305ep
```

will return a pretty printed table like this...

```text
Port ID |       Status | Power class | Voltage (V) | Current (mA) | Power (W) | Temperature (°C) | Error status
      1 |    Searching |             |           0 |            0 |  0.000000 |               32 | No Error
      2 |    Searching |             |           0 |            0 |  0.000000 |               33 | No Error
      3 |    Searching |             |           0 |            0 |  0.000000 |               33 | No Error
      4 |    Searching |             |           0 |            0 |  0.000000 |               33 | No Error

```
