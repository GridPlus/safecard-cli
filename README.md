# safecard-cli

A CLI for interacting with a GridPlus safecard through an HID card reader. Currently has one function, which exports a card's recovery seed. Call like so:

```
./safecard-cli exportSeed
```

# Build
```
make build
```

# Run development version
In development, the CLI can be run directly without first building a binary by running it like so:
```
go run main.go exportSeed
```