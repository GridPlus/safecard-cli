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
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// exportSeedCmd represents the exportSeed command
var exportSeedCmd = &cobra.Command{
	Use:   "exportSeed",
	Short: "Export our Safecard wallet's root seed",
	Long:  `Export your Safecard wallet's root seed.`,
	Run: func(cmd *cobra.Command, args []string) {
		exportSeed()
	},
}

func init() {
	rootCmd.AddCommand(exportSeedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportSeedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportSeedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func exportSeed() {
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
	seed, err := cs.ExportSeed()
	if err != nil {
		fmt.Println("unable to export seed. err: ", err)
	}
	fmt.Printf("recovery seed:\n0x%x\n", seed)
}
