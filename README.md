# IDÅSEN Desk - CLI


![License][license-badge]
![Go][go-version-badge]
![Version][release-version-badge]

## Why?

I own an IKEA IDÅSEN Desk, a standing desk equipped with Bluetooth functionality but lacking configurable buttons.
Unfortunately, the Bluetooth adapter only accommodates one connected device at a time. To overcome this limitation and
enable remote control from anywhere within or outside my house, I developed a command-line interface (CLI) tool. This
tool enables my Raspberry Pi to establish a connection with the desk and act as its controller. All my computers,
laptops, and phones are configured with SSH certificates, facilitating remote connections. Through SSH, these devices
can securely connect to the Raspberry Pi and control the desk.

### Caveat

While the underlying Bluetooth library supports Windows, Mac & Linux, the CLI only targets Linux since it was built 
directly on a Raspberry Pi. The underlying implementation for each operating system changes based on the build target.
This results in the build failing on anything other than Linux.

## Usage

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

To see any additional configuration options for a given command you can run the 
following to see the output for a single command.

```bash
desk help $COMMAND
# desk help stand
```

### Configure

The configure command will display a list of bluetooth devices currently
connected to your adapter or broadcasting. Selecting the given device will
save the localised name and address to the configuration file for execution.


<p>
    <img src="./assets/monitor_configure.gif" width="100%" alt="Desk Configuration">
</p>

### Standing

Move the desk from the current position to the configured standing position.

<p>
    <img src="./assets/monitor_stand.gif" width="100%" alt="Desk Standing">
</p>

### Sitting

Move the desk from the current position to the configured sitting position.

<p>
    <img src="./assets/monitor_sit.gif" width="100%" alt="Desk Sitting">
</p>

### Monitoring

Connect to the desk and monitor the height as it changes during manual operation.

<p>
    <img src="./assets/monitor_example.gif" width="100%" alt="Desk Monitoring">
</p>

[license-badge]: https://img.shields.io/github/license/stephensli/idasen-desk?style=flat-square

[go-version-badge]: https://img.shields.io/github/go-mod/go-version/stephensli/idasen-desk?style=flat-square

[release-version-badge]: https://img.shields.io/github/v/release/stephensli/idasen-desk?style=flat-square
