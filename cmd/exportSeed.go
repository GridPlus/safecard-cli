/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"
)

// exportSeedCmd represents the exportSeed command
var exportSeedCmd = &cobra.Command{
	Use:   "exportSeed",
	Short: "Export our Safecard wallet's root seed.",
	Long:  `Export your Safecard wallet's root seed.`,
	Run: func(cmd *cobra.Command, args []string) {
		exportSeed()
	},
}

func init() {
	rootCmd.AddCommand(exportSeedCmd)
}

func exportSeed() {
	seed, err := card.ExportSeed(readerIdx)
	if err != nil {
		return
	}
	fmt.Println("\n-------------------------")
	fmt.Println("Exporting wallet seed in hexadecimal format")
	fmt.Println("NOTE: This is NOT your seed phrase - it is the hash of your seed phrase. It will be limited in its use.")
	fmt.Println("If you wish to export private keys instead, please see the `exportPriv` command.")
	fmt.Println("-------------------------")
	fmt.Println()
	fmt.Printf("%x\n\n", seed)
}
