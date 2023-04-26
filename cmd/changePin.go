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

// changePinCmd represents the changePin command
var changePinCmd = &cobra.Command{
	Use:   "changePin",
	Short: "Changes the cards pin.",
	Long:  `Changes the cards pin.`,
	Run: func(cmd *cobra.Command, args []string) {
		changePin()
	},
}

func init() {
	rootCmd.AddCommand(changePinCmd)
}

func changePin() {

	//Open a secure connection
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
	fmt.Println("Please enter current 6 digit pin:")
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

	//Prompt user to confirm pin change before running
	confirm := promptui.Select{
		Label: "Yes/No",
		Items: []string{"Yes", "No"},
	}
	fmt.Println("Are you sure you want to change this cards pin?")
	_, result, err = confirm.Run()
	if err != nil || result != "Yes" {
		fmt.Println("aborting changePin command.")
		return
	}

	//Prompt user for new pin
	pin_promt := promptui.Prompt{
		Label: "New Pin",
		Mask:  '*',
	}
	fmt.Println("Please enter new 6 digit pin:")
	result, err = pin_promt.Run()
	if err != nil {
		fmt.Println("prompt failed: err: ", err)
		return
	}

	//Change pin
	err = cs.ChangePIN(result)
	if err != nil {
		fmt.Println("unable to change pin: err: ", err)
		return
	}
	fmt.Println("succesfully changed your pin!")
}
