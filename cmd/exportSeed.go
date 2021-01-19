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

	"github.com/GridPlus/keycard-go/gridplus"
	"github.com/gridplus/safecard-cli/card"
	"github.com/manifoldco/promptui"
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
	cs, err := card.OpenSecureConnection()
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
		fmt.Println("error verifying pin. err: ", err)
		return
	}
	seed, err := cs.ExportSeed()
	if err == gridplus.ErrSeedInvalidLength {
		fmt.Println("card does not appear to have valid exportable seed.")
		return
	}
	if err != nil {
		fmt.Println("unable to export seed. err: ", err)
		return
	}
	fmt.Printf("recovery seed:\n0x%x\n", seed)
}
