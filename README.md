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
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, -v             Enable verbose logging (default: false) [$VERBOSE]
   --stand-height value      The target end height for standing (default: 1.12)
   --sit-height value        The target end height for sitting (default: 0.74)
   --sit                     Put the desk into a sitting position (default: false)
   --stand                   Put the desk into a standing position (default: false)
   --target value, -t value  Move the desk into the target position (default: 0)
   --monitor, -m             Monitor the movement of the desk during manual movement (default: false)
   --help, -h                show help (default: false)
```

[license-badge]: https://img.shields.io/github/license/stephensli/idasen-desk?style=flat-square

[go-version-badge]: https://img.shields.io/github/go-mod/go-version/stephensli/idasen-desk?style=flat-square

[release-version-badge]: https://img.shields.io/github/v/release/stephensli/idasen-desk?style=flat-square
