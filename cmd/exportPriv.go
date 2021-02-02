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

	"github.com/gridplus/safecard-cli/card"
	"github.com/gridplus/safecard-cli/crypto"
	"github.com/spf13/cobra"
)

var coin string
var electrum bool
var electrumWallet bool
var numKeys int
var startPath string
var wif bool

// exportPriv represents the exportPriv command
var exportPriv = &cobra.Command{
	Use:   "exportPriv",
	Short: "Exports safecard private keys, with various options for BTC and ETH",
	Long: `Exports safecard private keys. Recommended usage below. Please see README for more details and options.

Ethereum:
The exported Ethereum private key(s) (printed as hexadecimal strings) may be pasted directly into MetaMask. By default, the Lattice only uses the first key, so you can simply run:

safecard-cli exportPriv --coin ETH

Bitcoin:
For Bitcoin we offer export of the master key for use in Electrum (recommended) as well as different types of individual account private keys.

(Recommended) Master Key (Electrum)
If you wish to import a full hierarchical deterministic (HD) wallet into Bitcoin wallet software, we highly recommend exporting the "master private key" and importing it into Electrum.

Note that the exported key is compatible with Electrum but probably not with anything else. Electrum expects a master key that is derived at the path m/49'/0'/0', whereas usually "master key" refers to an underived key.

safecard-cli exportPriv --electrum-master-priv
You can use the result of that to create an HD wallet in Electrum.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		exportPrivKey()
	},
}

func init() {
	rootCmd.AddCommand(exportPriv)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportPriv.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportPriv.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	exportPriv.Flags().BoolVar(&electrumWallet,
		"electrum-master-priv",
		false,
		"(Default false) If true, export a single BIP49 master private key which may be imported into Electrum")
	exportPriv.Flags().IntVar(&numKeys,
		"num-keys",
		1,
		"(Default 1) Number of private keys to export")
	exportPriv.Flags().StringVar(&startPath,
		"start-path",
		"m/49'/0'/0'/0/0",
		"(Default m/49'/0'/0'/0/0) A BIP32 path providing the starting point of the derivation")
	exportPriv.Flags().BoolVar(&electrum,
		"electrum",
		false,
		"(Default false) If true, includes prefixes and encoding necessary to import Bitcoin private keys into Electrum")
	exportPriv.Flags().StringVar(&coin,
		"coin",
		"",
		"(Default BTC) May be specified as either 'BTC' or 'ETH'. By default this is inferred from 'start-path' indices.")
	exportPriv.Flags().BoolVar(&wif,
		"wif",
		false,
		"(Default false) If true, single private keys are exported in WIF.")
}

func exportPrivKey() {
	seed, err := card.ExportSeed()
	if err != nil {
		return
	}

	// Export individual private key(s)
	if false == electrumWallet {
		path, err := crypto.GetPath(startPath)
		if err != nil {
			fmt.Println("Error encountered parsing path: ", err)
			return
		}
		// If the user specified a currency we can change the path
		if coin == "ETH" {
			path[0] = 44 + 0x80000000
			path[1] = 60 + 0x80000000
		}
		// Update the path string for display to capture any changes that might have happened
		startPath = crypto.GetStrPath(path)
		// Determine the currency to use based on the path
		currency, err := crypto.GetCurrencyType(path)
		if err != nil {
			fmt.Println("Error encountered with path: ", err)
		}
		// Set wif=true if we are exporting keys for Electrum
		if true == electrum {
			wif = true
		}

		fmt.Println("\n-------------------------")
		fmt.Printf("Exporting %d %s private keys\n", numKeys, currency)
		fmt.Printf("Start path: %s\n", startPath)
		fmt.Println("-------------------------")
		fmt.Println()
		for i := 0; i < numKeys; i++ {
			priv, err := crypto.DerivePrivateKey(seed, path, wif)
			if err != nil {
				fmt.Println("Error encountered deriving key: ", err)
				return
			}
			prefix := ""
			if true == electrum {
				prefix, err = crypto.GetElectrumPrivKeyPrefix(path)
				if err != nil {
					fmt.Println("Error encountered buidling Electrum import string: ", err)
				}
			}
			fmt.Printf("%s%s\n", prefix, priv)
			path[len(path)-1]++
		}
		fmt.Println()

		return
	}

	// Otherwise we are exporting the Electrum master private key
	fmt.Println("\n-------------------------")
	fmt.Println("Exporting root private key for Electrum (m/49'/0'/0')")
	fmt.Println("You can use this as your 'master private key' when creating an Electrum wallet.")
	fmt.Println("-------------------------")
	fmt.Println()
	privateKey, err := crypto.GetElectrumBIP49MasterPriv(seed)
	if err != nil {
		return
	}
	fmt.Printf("%s\n\n", privateKey)
}
