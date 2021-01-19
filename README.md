# safecard-cli

A CLI for interacting with a GridPlus safecard through an HID card reader. Currently has one function, which exports a card's recovery seed. Call like so:

```
./safecard-cli exportSeed
```

## Build
```
make build
```

## Run development version
In development, the CLI can be run directly without first building a binary by running it like so:
```
go run main.go exportSeed
```

## Development

### Adding a new CLI command
In order to develop a new command for the CLI (e.g. exportSeed or deleteSeed) one should use the cobra autogenerate tool to set up a preformatted file under the cmd/ directory, by using the command below.
```
cobra add $commandName
```

This will autogenerate the necessary file under the cmd/ directory for the new shell command
Further details on the cobra generator here: https://github.com/spf13/cobra/blob/master/cobra/README.md