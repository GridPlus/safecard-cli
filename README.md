# safecard-cli

A CLI for interacting with a GridPlus safecard through an HID card reader. Runs on both Windows and MacOS with a USB HID card reader attached. Here's an example of an [HID USB card reader](https://www.amazon.com/HID-OMNIKEY-3121-Card-Reader/dp/B00AT4NX8S/ref=sr_1_14?dchild=1&keywords=hid+reader&qid=1611873802&sr=8-14)

## Usage

Download the appropriate binary for your platform from the Release page here: https://github.com/GridPlus/safecard-cli/releases/tag/latest

| OS    | Binary |
|:------|:-------|
| MacOS | safecard-cli |
| Windows | safecard-cli |

Run the binary on it's own to see usage info
```
safecard-cli
```

There are currently two implemented command, exportSeed and deleteSeed, each requires the user pin to be verified before executing.

### Export Seed

```
safecard-cli exportSeed
```

Export seed will export the card's root wallet key as a binary seed represented in hex. This hex seed can be used to derive wallet private keys and addresses.

An example means of [using the exported seed](#using-the-exported-seed) is shown below.

### Delete Seed

```
safecard-cli deleteSeed
```

Delete seed will delete the master wallet key, effectively destroying the safecard's wallet. This operation is irreversible and requires a pin and confirmation from the user before the wallet is able to be deleted.

### Using the Exported Seed
A simple NodeJS class that you could use to derive private keys and addresses is shown below. You can import a single private key into metamask for Ethereum. For Bitcoin, you would need to use something like bitcoinjs-lib (https://www.npmjs.com/package/bitcoinjs-lib) to build the transactions.

```
const bip32 = require('bip32');
const bitcoin = require('bitcoinjs-lib');
const ethereum = require('ethereumjs-util');
const ETH_PATH = "m/44'/60'/0'/0"
const BTC_PATH = "m/49'/0'/0'/0"
class Wallet {
  constructor(seed) {
    if (!seed)
      throw new Error('You must include your seed as a hex string');
    this.wallet = bip32.fromSeed(Buffer.from(seed, 'hex'))
  }
  getAddress(type, idx) {
    let priv = this.getPriv(type, idx);
    if (type === 'ETH') {
      return ethereum.privateToAddress(priv).toString('hex');
    } else if (type === 'BTC') {
      const data = bitcoin.payments.p2sh({
        redeem: bitcoin.payments.p2wpkh({ pubkey: pub }),
      });
      return data.address;
    } else {
      throw new Error(`Unsupported derivation type: ${type}`);
    }
  }
  getPriv(type, idx) {
    if (type === 'ETH') {
      return this._getPriv(`${ETH_PATH}/${idx}`)
    } else if (type === 'BTC') {
      return this._getPriv(`${BTC_PATH}/${idx}`)
    } else {
      throw new Error(`Unsupported derivation type: ${type}`);
    }
  }
  _getPriv(path) {
    return this.wallet.derivePath(path).privateKey;
  }
}
```
## Build
For Mac
```
make build
```
For Windows
```
make windows-build
```

## Development

### Run development version
In development, the CLI can be run directly without first building a binary by running it like so:
```
go run main.go exportSeed
```
### Adding a new CLI command
In order to develop a new command for the CLI (e.g. exportSeed or deleteSeed) one should use the cobra autogenerate tool to set up a preformatted file under the cmd/ directory, by using the command below.
```
cobra add $commandName
```

This will autogenerate the necessary file under the cmd/ directory for the new shell command
Further details on the cobra generator here: https://github.com/spf13/cobra/blob/master/cobra/README.md