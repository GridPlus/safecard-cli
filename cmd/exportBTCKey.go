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

var electrumPrefixes bool
var electrumWallet bool
var numKeys int
var startPath string

// exportBTCKeyCmd represents the exportBTCKey command
var exportBTCKeyCmd = &cobra.Command{
	Use:   "exportBTCKey",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		exportBTCKey()
	},
}

func init() {
	rootCmd.AddCommand(exportBTCKeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportBTCKeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportBTCKeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	exportBTCKeyCmd.Flags().BoolVar(&electrumWallet,
		"electrum-master-priv",
		false,
		"(Default false) If true, export a single BIP49 master private key which may be imported into Electrum")
	exportBTCKeyCmd.Flags().IntVar(&numKeys,
		"num-keys",
		1,
		"(Default 1) Number of private keys to export")
	exportBTCKeyCmd.Flags().StringVar(&startPath,
		"start-path",
		"m/49'/0'/0'/0/0",
		"(Default m/49'/0'/0'/0/0) A BIP32 path providing the starting point of the derivation")
	exportBTCKeyCmd.Flags().BoolVar(&electrumPrefixes,
		"electrum-prefixes",
		false,
		"(Default false) If true, include prefixes on exported individual private keys that allows them to be properly imported into Electrum.")
}

func exportBTCKey() {
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
		currency, err := crypto.GetCurrencyType(path)
		if err != nil {
			fmt.Println("Error encountered with path: ", err)
		}
		fmt.Println("\n-------------------------")
		fmt.Printf("Exporting %d %s private keys\n", numKeys, currency)
		fmt.Printf("Start path: %s\n", startPath)
		fmt.Println("-------------------------")
		fmt.Println()
		for i := 0; i < numKeys; i++ {
			priv, err := crypto.DerivePrivateKey(seed, path)
			if err != nil {
				fmt.Println("Error encountered deriving key: ", err)
				return
			}
			prefix, err := crypto.GetElectrumPrivKeyPrefix(path)
			if err != nil {
				fmt.Println("Error encountered buidling Electrum import string: ", err)
			}
			fmt.Printf("%s%s\n", prefix, priv)
			path[len(path)-1]++
		}
		fmt.Println()

		return
	}

	// Otherwise we are exporting the Electrum master private key
	privateKey, err := crypto.GetElectrumBIP49MasterPriv(seed)
	if err != nil {
		return
	}
	fmt.Println("private key: ")
	fmt.Println(privateKey)
}
