# IDÅSEN Desk - CLI

![License][license-badge]
![Go][go-version-badge]
![Version][release-version-badge]

```bash
NAME:
   Idasen CLI - A simple CLI to interface with the Idasen desk

USAGE:
   desk-cli [global options] command [command options] [arguments...]

COMMANDS:
   configure  configure the device to connect to.
   stand      Move the desk to the configured standing position.
   sit        Move the desk to the configured sitting position.
   position   Move the desk to the provided position value.
   toggle     Toggle the desk height between standing and sitting.
   monitor    Monitor and log the position of the desk as it moves
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

[license-badge]: https://img.shields.io/github/license/stephensli/idasen-desk?style=flat-square

[go-version-badge]: https://img.shields.io/github/go-mod/go-version/stephensli/idasen-desk?style=flat-square

[release-version-badge]: https://img.shields.io/github/v/release/stephensli/idasen-desk?style=flat-square
