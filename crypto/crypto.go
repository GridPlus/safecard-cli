package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/FactomProject/basen"
	log "github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip32"
)

// B58Enc is the base58 library used for base58 encoding data
var B58Enc = basen.NewEncoding("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// The BIP49 "version" is different than the normal BIP44 one. We need this in order
// to produce the `yprv` prefixed master key (as oposed to xprv for BIP44)
// see: https://github.com/iancoleman/bip39/commit/93c3ef47579733040dbc6eec865b528d1ca49911#diff-e80f5b5b0593f29b12532598b7f4264308ce02c03e85bd8bc3e4e5d9cb5b3a90R265
func bip49Version() []byte {
	BIP49Version, _ := hex.DecodeString("049d7878")
	return BIP49Version
}

// Convenience function for generating a 4-byte checksum from a preimage
func doubleSha256Checksum(data []byte) ([]byte, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, err
	}
	hash1 := hasher.Sum(nil)
	hasher = sha256.New()
	_, err = hasher.Write(hash1)
	if err != nil {
		return nil, err
	}
	hash2 := hasher.Sum(nil)
	return hash2[:4], nil
}

// Get the WIF representation of a private key. This is used in Electrum to import single account keys.
// See: https://en.bitcoin.it/wiki/Wallet_import_format
func getWif(priv []byte) (string, error) {
	version := []byte{0x80}
	compression := []byte{0x01}
	key := append(version, priv...)
	key = append(key, compression...)
	// Add checksum
	cs, err := doubleSha256Checksum(key)
	if err != nil {
		return "", err
	}
	key = append(key, cs...)
	// Convert to base58 string and return
	return B58Enc.EncodeToString(key), nil
}

// DerivePrivateKey derives a single HD wallet's private key given the seed and a path of indices.
// Returns hex string representation of private key
func DerivePrivateKey(seed []byte, path []uint32, wif bool) (privKey string, err error) {
	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Error("error deriving master private key from seed: ", err)
		return "", err
	}
	for i := 0; i < len(path); i++ {
		key, err = key.NewChildKey(path[i])
		if err != nil {
			return "", err
		}
	}
	if true == wif {
		return getWif(key.Key)
	}
	keyStr := hex.EncodeToString(key.Key)
	return keyStr, nil
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
		return "", err
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

// GetElectrumPrivKeyPrefix returns a string containing a prefix that must be included
// when a user imports a private key string into Electrum (Bitcoin wallet).
// The Lattice uses BIP49 derivations, so we expect 49' most of the time, though 44'
// is also supported (legacy p2pkh, so there is no prefix needed).
// Otherwise an error is thrown.
func GetElectrumPrivKeyPrefix(path []uint32) (prefix string, err error) {
	switch path[0] {
	case 49 + 0x80000000:
		return "p2wpkh-p2sh:", nil // Wrapped segwit (default)
	case 44 + 0x80000000:
		return "", nil // Default p2pkh
	default:
		return "", errors.New("unsupported path type. Must be 49' or 44'")
	}
}

// GetPath takes a path string (e.g. m/44'/0'/0'/0/0) and converts it to uint32 indices
func GetPath(pathStr string) ([]uint32, error) {
	indices := strings.Split(pathStr, "/")
	start := 0
	if indices[0] == "m" {
		start = 1
	}
	path := make([]uint32, (len(indices) - start))
	for i := start; i < len(indices); i++ {
		_strIdx := indices[i]
		_L := len(indices[i])
		isHardened := string(indices[i][(_L-1):]) == "'"
		strIdx := _strIdx
		if isHardened {
			strIdx = _strIdx[:(_L - 1)]
		}
		idx, err := strconv.ParseUint(strIdx, 10, 32)
		if err != nil {
			return nil, err
		}
		if isHardened {
			idx += 0x80000000
		}
		path[i-start] = uint32(idx)
	}
	return path, nil
}

// GetStrPath converts a set of indices to a BIP32 path string (e.g. m/44'/0'/0'/0/0)
func GetStrPath(path []uint32) string {
	pathStr := "m"
	for i := 0; i < len(path); i++ {
		pathStr += "/"
		if path[i] >= uint32(0x80000000) {
			pathStr += strconv.FormatUint(uint64(path[i])-0x80000000, 10)
			pathStr += "'"
		} else {
			pathStr += strconv.FormatUint(uint64(path[i]), 10)
		}
	}
	return pathStr
}

// GetCurrencyType takes a path (set of u32 indices) and returns the name of the
// currency, if applicable. Returns error if currency is not supported or path
// is malformed.
func GetCurrencyType(path []uint32) (name string, err error) {
	if len(path) < 2 {
		return "", errors.New("path contains <2 indices")
	}
	switch path[1] {
	case 0x80000000:
		return "Bitcoin", nil
	case 60 + 0x80000000:
		return "Ethereum", nil
	default:
		return "", errors.New("currency must be either 0' (Bitcoin) or 60' (Ethereum)")
	}
}
