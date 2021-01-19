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
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// deleteSeedCmd represents the deleteSeed command
var deleteSeedCmd = &cobra.Command{
	Use:   "deleteSeed",
	Short: "WARNING: This command is irreversible once completed. Deletes the safecard root seed",
	Long: `WARNING: This command is irreversible once completed.

	Deletes the safecard's root seed, which is the root of all of the card's wallet keys.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteSeed()
	},
}

func init() {
	rootCmd.AddCommand(deleteSeedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteSeedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteSeedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deleteSeed() {
	cs, err := card.Connect()
	if err != nil {
		fmt.Println("error connecting to card")
		fmt.Println(err)
		return
	}
	err = cs.Select()
	if err != nil {
		fmt.Println("error selecting applet. err: ", err)
		return
	}
	err = cs.Pair()
	if err != nil {
		fmt.Println("error pairing with card. err: ", err)
		return
	}
	err = cs.OpenSecureChannel()
	if err != nil {
		fmt.Println("error opening secure channel. err: ", err)
		return
	}
	//Prompt user for pin
	prompt := promptui.Prompt{
		Label: "Pin",
		Mask:  '*',
	}
	fmt.Println("Please enter 6 digit pin")
	result, err := prompt.Run()
	if err != nil {
		fmt.Println("prompt failed: err: ", err)
		return
	}

	err = cs.VerifyPIN(result)
	if err != nil {
		fmt.Println("error verifying pin. err: ", err)
		return
	}
	//Prompt user to confirm seed deletion before running
	confirm := promptui.Select{
		Label: "Yes/No",
		Items: []string{"yes", "no"},
	}
	fmt.Println("Are you sure you want to delete this card's master seed? This action is irreversible")
	_, result, err = confirm.Run()
	if err != nil || result != "yes" {
		fmt.Println("aborting deleteSeed command")
		return
	}
	err = cs.RemoveKey()
	if err != nil {
		fmt.Println("unable to remove master key: ", err)
		return
	}
	fmt.Println("succesfully removed master key!")
}
