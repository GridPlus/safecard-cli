package crypto

import (
	"encoding/hex"
	"fmt"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	log "github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip32"
)

const ETHDerivationPath = "m/44'/60'/0'/0/0"
const BTCDerivationPath = "m/44'/0'/0'/0/0"

// The BIP49 "version" is different than the normal BIP44 one. We need this in order
// to produce the `yprv` prefixed master key (as oposed to xprv for BIP44)
// see: https://github.com/iancoleman/bip39/commit/93c3ef47579733040dbc6eec865b528d1ca49911#diff-e80f5b5b0593f29b12532598b7f4264308ce02c03e85bd8bc3e4e5d9cb5b3a90R265
func bip49Version() []byte {
	BIP49Version, _ := hex.DecodeString("049d7878")
	return BIP49Version
}

//TODO: Clean up
func DeriveWalletPrivateKey(seed []byte) (address string, privKey string, err error) {
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		log.Error("could not parse seed into wallet. ", err)
		return "", "", err
	}
	path := hdwallet.MustParseDerivationPath(ETHDerivationPath)
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Error("could not derive path: ", err)
		return "", "", err
	}

	fmt.Println(account.Address.Hex())
	privateKey, err := wallet.PrivateKeyHex(account)
	if err != nil {
		log.Error("could not derive private key from account: ", err)
		return "", "", err
	}

	return account.Address.Hex(), privateKey, nil
}

// GetElectrumBIP49MasterPriv exports a BIP49 master private key for import into Electrum.
// Generally, the master private key would be the `masterKey` defined below.
// However, Electrum requires further derivation: m/49'/0'/0'. They consider this derived
// key to be the "master private key". As far as I can tell this breaks from other wallet
// software, but Electrum is the only popular one that lets you import a master private key
// for an HD wallet.
func GetElectrumBIP49MasterPriv(seed []byte) (privKey string, err error) {
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Error("error deriving master private key from seed: ", err)
	}
	// Replace the key version
	masterKey.Version = bip49Version()
	// Derive m/49'/0'/0' to get to the key which Electrum will consider the "master"
	hardened := uint32(0x80000000)
	first, err := masterKey.NewChildKey(49 + hardened)
	if err != nil {
		fmt.Println("err deriving child key")
		return "", err
	}
	second, err := first.NewChildKey(hardened)
	if err != nil {
		fmt.Println("err deriving child key")
		return "", err
	}
	third, err := second.NewChildKey(hardened)
	if err != nil {
		fmt.Println("err deriving child key")
		return "", err
	}
	third.Version = bip49Version()
	return third.B58Serialize(), nil
}
