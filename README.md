# safecard-cli

A CLI for interacting with a GridPlus safecard through an HID card reader. Runs on both Windows and MacOS with a USB HID card reader attached. Here's an example of an [HID USB card reader](https://www.amazon.com/HID-OMNIKEY-3121-Card-Reader/dp/B00AT4NX8S/ref=sr_1_14?dchild=1&keywords=hid+reader&qid=1611873802&sr=8-14)

## Usage

Download the appropriate binary for your platform from the Release page here: https://github.com/GridPlus/safecard-cli/releases/tag/latest

| OS    | Binary |
|:------|:-------|
| MacOS | safecard-cli |
| Windows | safecard-cli.exe |

Run the binary on its own to see usage info
```
safecard-cli
```

There are currently two implemented commands, exportSeed and deleteSeed.
Both require the user pin to be verified before executing.

### Export Seed

```
safecard-cli exportSeed
```

Export seed will export the card's master wallet seed as a binary seed represented in hex. This hex seed can be used to derive wallet private keys and addresses.

An example means of [using the exported seed](#using-the-exported-seed) is shown below.

### Delete Seed

```
safecard-cli deleteSeed
```

Delete seed will delete the master wallet key, effectively destroying the safecard's wallet. This operation is irreversible and requires a pin and confirmation from the user before the wallet is able to be deleted.

### Using the Exported Seed
SafeCards store the entropy (i.e. “seed”) of a hierarchical, deterministic BIP32 wallet. This is notably different from a seed phrase, specifically it is the hash of a phrase (plus optional password). This means that you cannot “go back” to a seed phrase, but your seed is all you need to derive addresses and private keys.

You can use any number of off-the-shelf developer tools to derive keys from a seed. When you get the desired private key, you can either import it directly into MetaMask (for use with Ethereum) or use it to sign transactions you build using a Bitcoin library like [bitcoinjs-lib ](https://www.npmjs.com/package/bitcoinjs-lib).

The following script is an example of how to take your seed and derive addresses or private keys using Javascript modules. If you wish to do this at the command line, we recommend putting the following script in its own directory (name it derive.js or something) and running npm init, pressing enter a bunch of times until the prompt is gone, and then running npm i --save bip32 bitcoinjs-lib ethereumjs-util.  Then you can derive your keys and addresses like this:
node derive.js 6cc741dc06b353b97852b15c42d0fcb672d48983630840d13780715dd23f6655e7344ff000122e078e7f6b82edb7d4225f15767c61f3ab9a1400ce9f42d38cd9 1 ETH priv

```
const bip32 = require('bip32');
const bitcoin = require('bitcoinjs-lib');
const ethereum = require('ethereumjs-util');
const ETH_PATH = "m/44'/60'/0'/0"
const BTC_PATH = "m/49'/0'/0'/0"
const seed = process.argv[2];
if (!seed)
  throw new Error('You must include your seed as a hex string');
const wallet = bip32.fromSeed(Buffer.from(seed, 'hex'));
let idx = process.argv[3];
if (idx === undefined || isNaN(idx)) {
  idx = 0;
  console.warn('Unspecified derivation index. Using default index 0.');
}
let type = process.argv[4];
if (type.toUpperCase() !== 'ETH' && type.toUppderCase() !== 'BTC') {
  type = 'ETH';
  console.warn('Unspecified derivation type. Using default type ETH.');
}
let showPriv = false;
if (process.argv[5] === 'priv') {
  showPriv = true;
}
const path = type === 'ETH' ? ETH_PATH : BTC_PATH;
const key = wallet.derivePath(`${path}/${idx}`);
const priv = key.privateKey;
let addr;
if (type === 'ETH') {
  addr = '0x' + ethereum.privateToAddress(priv).toString('hex')
} else {
  const preAddr = bitcoin.payments.p2sh({
    redeem: bitcoin.payments.p2wpkh({ pubkey: key.publicKey }),
  });
  addr = preAddr.address;
}
console.log('\n\n')
console.log(`---${type} key index ${idx}---`)
console.log(`Address: ${addr}`);
if (showPriv === true) {
  console.log(`Private Key: ${priv.toString('hex')}`)
}
console.log('------------------------'
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