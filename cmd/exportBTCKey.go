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
}

func exportBTCKey() {
	seed, err := card.ExportSeed()
	if err != nil {
		return
	}

	privateKey, err := crypto.GetElectrumBIP49MasterPriv(seed)
	if err != nil {
		return
	}
	fmt.Println("private key: ")
	fmt.Println(privateKey)
}
