package crypto

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	log "github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip32"
)

const ETHDerivationPath = "m/44'/60'/0'/0/0"
const BTCDerivationPath = "m/44'/0'/0'/0/0"

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

func DeriveBTCPrivateKey(seed []byte, index int) (address []byte, privKey string, err error) {

	// wallet, err := hdwallet.NewFromSeed(seed)
	// if err != nil {
	// 	log.Error("could not parse seed into wallet. ", err)
	// 	return "", "", err
	// }
	// path := hdwallet.MustParseDerivationPath(ETHDerivationPath)
	// account, err := wallet.Derive(path, false)
	// if err != nil {
	// 	log.Error("could not derive path: ", err)
	// 	return "", "", err
	// }

	//BTC Wallet Implementation

	// masterPrivKey := btcwallet.MasterKey(seed[0:32])

	//BIP32 attempt
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Error("error deriving master private key from seed: ", err)
	}

	// fmt.Println("master private key")
	// fmt.Println(masterPrivKey.B58Serialize())
	serialized, _ := masterKey.Serialize()
	fmt.Println("master key: ")
	fmt.Printf("%X\n", serialized)
	base58Key := masterKey.B58Serialize()
	fmt.Println("base58 master key")
	fmt.Println(base58Key)

	hardened := uint32(0x80000000)
	first, err := masterKey.NewChildKey(49 + hardened)
	if err != nil {
		fmt.Println("err deriving child key")
		return
	}
	second, err := first.NewChildKey(hardened)
	if err != nil {
		fmt.Println("err deriving child key")
		return
	}
	third, err := second.NewChildKey(hardened)
	if err != nil {
		fmt.Println("err deriving child key")
		return
	}
	fourth, err := third.NewChildKey(0)
	if err != nil {
		fmt.Println("err deriving child key")
		return
	}
	fifth, err := fourth.NewChildKey(0)
	if err != nil {
		fmt.Println("err deriving child key")
		return
	}

	fmt.Println("raw privKey bytes")
	fmt.Printf("%X\n", fifth.Key)

	fifthPrivKey := SerializeShort(fifth)
	fmt.Println("fifth private key")
	fmt.Printf("%X\n", fifthPrivKey)
	return nil, fifth.String(), nil

	// address = fifth.Address()

	// return address, fifth.String(), nil
	// address = privKeyWallet.Address()

	// return address, privKeyWallet.String(), nil

	//bip32 implementation
	// masterKey, err := bip32.NewMasterKey(seed)
	// if err != nil {
	// 	log.Error("could not parse seed into master key: ", err)
	// 	return "", "", err
	// }
	// return

	// privKey, err = masterKey.NewChildKey(index)
	// if err != nil {
	// 	log.Error("could not derive child key from master: ", err)
	// 	return "", "", err
	// }
	//Derive Public Key
	//Derive Address
}

func NewUnsaltedMasterKey(seed []byte) (*bip32.Key, error) {
	// Generate key and chaincode
	hmac := hmac.New(sha512.New, nil)
	_, err := hmac.Write(seed)
	if err != nil {
		return nil, err
	}
	intermediary := hmac.Sum(nil)

	// Split it into our key and chain code
	keyBytes := intermediary[:32]
	chainCode := intermediary[32:]

	// // Validate key
	// err = validatePrivateKey(keyBytes)
	// if err != nil {
	// 	return nil, err
	// }

	// PrivateWalletVersion is the version flag for serialized private keys
	PrivateWalletVersion, _ := hex.DecodeString("0488ADE4")
	// Create the key struct
	key := &bip32.Key{
		Version:     PrivateWalletVersion,
		ChainCode:   chainCode,
		Key:         keyBytes,
		Depth:       0x0,
		ChildNumber: []byte{0x00, 0x00, 0x00, 0x00},
		FingerPrint: []byte{0x00, 0x00, 0x00, 0x00},
		IsPrivate:   true,
	}

	return key, nil
}

func SerializeShort(key *bip32.Key) []byte {
	// Private keys should be prepended with a single null byte
	keyBytes := key.Key
	if key.IsPrivate {
		keyBytes = append([]byte{0x0}, keyBytes...)
	}
	return keyBytes
}
