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
	"strings"

	"github.com/gridplus/safecard-cli/card"
	"github.com/gridplus/safecard-cli/util"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// exportMnemonicCmd represents the exportData command
var exportMnemonicCmd = &cobra.Command{
	Use:   "exportMnemonic",
	Short: "Export our Safecard wallet's mnemonic phrase.",
	Long:  `Export your Safecard wallet's mnemonic phrase.`,
	Run: func(cmd *cobra.Command, args []string) {
		exportMnemonic()
	},
}

func init() {
	rootCmd.AddCommand(exportMnemonicCmd)
}

func exportMnemonic() {
	cs, err := card.OpenSecureConnection(readerIdx)
	if err != nil {
		fmt.Println("unable to open secure connection with card: ", err)
		return
	}

	//Prompt user for pin
	prompt := promptui.Prompt{
		Label: "Pin",
		Mask:  '*',
	}
	fmt.Println("Please enter 6 digit pin:")
	result, err := prompt.Run()
	if err != nil {
		fmt.Println("prompt failed: err: ", err)
		return
	}
	err = cs.VerifyPIN(result)
	if err != nil {
		fmt.Println("Error validating PIN.")
		return
	}

	// Make sure the user is comfortable printing the mnemonic
	confirm := promptui.Select{
		Label: "Yes/No",
		Items: []string{"Yes", "No"},
	}
	fmt.Println(
		"Are you sure you want to print your mnemonic phrase?\n" +
		"Please make sure you are on a secure terminal and no one is looking over your shoulder.")
	_, result, err = confirm.Run()
	if err != nil || result != "Yes" {
		fmt.Println("Aborting")
		return
	}

	// Get the data region
	data, err := cs.ExportData()
	if err != nil {
		fmt.Println("error exporting data. err: ", err)
		return
	}
	// Parse out the mnemonic
	mnemonic := util.GetExportDataItemByTag(util.ExportDataTagMnemonic, data)
	if mnemonic == nil {
		fmt.Println(
			`No mnemonic found on the card. SafeCard version may be too low. 
			Only newer cards (v2.3+) support this feature.`)
		return
	}

	words, err := util.ConvertMnemonicIndicesToWords(mnemonic)
	fmt.Println("\nYour mnemonic phrase is:")
	fmt.Println("-------------------------")
	fmt.Println(strings.Join(words[:], " "))
	fmt.Println("-------------------------")
}

