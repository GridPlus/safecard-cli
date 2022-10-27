/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"reflect"

	"github.com/GridPlus/keycard-go/gridplus"

	"github.com/gridplus/safecard-cli/card"
	"github.com/spf13/cobra"
)

// checkCertCmd fetches and reports the cert on the card
var checkCertCmd = &cobra.Command{
	Use:   "checkCert",
	Short: "Verify the authenticity of a GridPlus SafeCard",
	Long: `Export and verify the certificate to ensure the card
				was programmed and validated by GridPlus.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkCert()
	},
}

func init() {
	rootCmd.AddCommand(checkCertCmd)
}

func checkCert() {
	// 1. Connect to the card
	cs, err := card.Connect(readerIdx)
	if err != nil {
		fmt.Println("ERROR: connecting to card", err)
		return
	}
	err = cs.Select()
	if err != nil {
		fmt.Println("ERROR: failed to select", err)
	}
	// 2. Export the certificate from the card
	cert, err := cs.ExportCert()
	if err != nil {
		fmt.Println("ERROR: failed to export cert", err)
		return
	}
	if len(cert) == 0 {
		fmt.Println("ERROR: Cert export not supported by your card. >=V2.4 SafeCard required.")
		return
	}
	fmt.Println("cert", cert)
	off := 0
	if cert[off] != 0x30 {
		fmt.Println("ERROR: bad certificate data")
		return
	}
	off++
	certLen := cert[off]
	off++
	// Slice out permissions (not currently used/enforeced)
	permissions := cert[off:off+4]
	off += 4
	// Slice out putkey
	pubkey := cert[off:off+2+65]
	off += 2 + 65
	// Slice out the signature
	sig := cert[off:certLen]

	// 3. Validate the cert - this verifies the card's ID pubkey was
	// signed by the GridPlus signing key
	fullCert := gridplus.SafecardCert{
		Permissions: permissions,
		PubKey: pubkey,
		Sig: sig,
	}
	valid := gridplus.ValidateCardCertificate(fullCert)
	if !valid {
		fmt.Println("ERROR: Certificate invalid")
		return
	}

	// 4. Ask for a signature from the card and confirm it was signed
	// by the correct ID pubkey
	challengeHash := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
		9, 10, 11, 12, 13, 14, 15, 16,
		17, 18, 19, 20, 21, 22, 23, 24,
		25, 26, 27, 28, 29, 30, 31, 32,
	}
	id, err := cs.IdentifyCard(challengeHash)
	if err != nil {
		fmt.Println("ERROR: failed to identify card", err)
		return
	}
	idOff := 0
	if id[idOff] != 0x80 {
		fmt.Println("ERROR: idenfity card returned bad data")
		return
	}
	idOff++
	idPubLen := int(id[idOff])
	idOff++
	idPubPlusMetadata := id[idOff-2:idOff+idPubLen]
	idPub := id[idOff:idOff+idPubLen]
	idOff += idPubLen
	idSig := id[idOff:]
	// Make sure the pubkey matches
	if !reflect.DeepEqual(idPubPlusMetadata, pubkey) {
		fmt.Println("ERROR: ID pubkey does not match cert!")
		return
	}
	challengeValid := gridplus.ValidateECDSASignature(idSig, idPub, challengeHash)
	if !challengeValid {
		fmt.Println("ERROR: Failed to validate ID key challenge")
		return
	}
	fmt.Println("SafeCard is valid and has been certified by GridPlus!")
}

