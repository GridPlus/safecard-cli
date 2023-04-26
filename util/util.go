package util

import (
	"encoding/binary"
	"errors"
)

const (
	ExportDataTagMnemonic = 0x01
)

/**
Search for a particular type of data in the full export data returned by the card
@param {tag} uint8 - the tag of the data to search for
@param {data} []byte - the full export data returned by the card
@return {[]byte} - the data item with the given tag, or nil if not found
*/
func GetExportDataItemByTag(tag uint8, data []byte) []byte {
	for i := 0; i < len(data); i++ {
		if i + 2 >= len(data) {
			return nil
		}
		itemTag := int(data[i])
		itemSz := int(data[i + 1])
		if i + 2 + itemSz >= len(data) {
			return nil
		}
		if itemTag == int(tag) {
			return data[(i + 2):(i + 2 + itemSz)]
		}
		i += (2 + itemSz)
	}
	return nil
}

func ConvertMnemonicIndicesToWords(data []byte) ([]string, error) {
	// Each word is represented by a u16 index corresponding to the word
	// in BIP39_WORD_LIST
	words := make([]string, len(data) / 2)
	for i := 0; i < len(data); i+=2 {
		wordIdx := binary.LittleEndian.Uint16(data[i:(i+2)])
		if wordIdx >= uint16(len(BIP39_WORD_LIST)) {
			return nil, errors.New("invalid BIP39 word index")
		}
		words[i / 2] = BIP39_WORD_LIST[wordIdx]
	}
	return words, nil
}