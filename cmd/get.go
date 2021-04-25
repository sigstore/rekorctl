// Copyright 2021 The Sigstore Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io/ioutil"

	"github.com/google/trillian"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type getProofResponse struct {
	Status string
	Proof  *trillian.GetInclusionProofByHashResponse
	Key    []byte
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Rekor get command",
	Long: `Performs a proof verification that a file
exists within the transparency log`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rekorServer := viper.GetString("rekor_server")
		url := rekorServer + "/api/v1/getproof"
		rekord := viper.GetString("rekord")

		rekorEntry, err := ioutil.ReadFile(rekord)
		if err != nil {
			return err
		}
		err = DoGet(url, rekorEntry)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
