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

## Development a New CLI Command
Using the cobra framework, one can add a new command for development with the following command:
```
cobra add $commandName
```

Further details on the cobra generator here: https://github.com/spf13/cobra/blob/master/cobra/README.md