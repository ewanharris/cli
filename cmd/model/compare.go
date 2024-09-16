/*
Copyright © 2024 OpenFGA

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

package model

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/openfga/cli/internal/authorizationmodel"
)

func getModel(cmd *cobra.Command, flagName string) (*authorizationmodel.AuthzModel, error) {
	authModel, err := authorizationmodel.Read(cmd, flagName)
	if err != nil {
		return nil, fmt.Errorf("failed to read model: %w", err)
	}
	return authModel, nil
}

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Perform a comparison between two models.",
	Long:  "Models are converted to the DSL format so that the ordering is matched and comments are stripped. The results DSLs are then compared.", //nolint:lll
	RunE: func(cmd *cobra.Command, _ []string) error {
		modelA, err := getModel(cmd, "model-a")
		if err != nil {
			return fmt.Errorf("failed to get model-a dsl: %w", err)
		}

		modelB, err := getModel(cmd, "model-b")
		if err != nil {
			return fmt.Errorf("failed to get model-a dsl: %w", err)
		}

		quiet, err := cmd.Flags().GetBool("quiet")
		if err != nil {
			return fmt.Errorf("failed to get quiet flag due to %w", err)
		}

		diff, err := modelA.DiffDSL(modelB)
		if err != nil {
			return fmt.Errorf("failed to compare DSL %w", err)
		}

		if diff == "" {
			return nil
		}

		errorString := "Models are not equal"
		if !quiet {
			errorString += "\n" + diff
		}

		return errors.New(errorString) //nolint:err113
	},
}

func init() {
	compareCmd.Flags().String("model-a", "latest", "The base model in the comparison")
	compareCmd.Flags().String("model-b", "", "The model to compare against")
	compareCmd.Flags().String("store-id", "", "Store ID to read from if any model is remote")
	compareCmd.Flags().BoolP("quiet", "q", false, "Do not print a diff of the model, only an error message,")

	if err := compareCmd.MarkFlagRequired("model-a"); err != nil {
		fmt.Printf("error setting flag required %v", err)
		os.Exit(1)
	}

	if err := compareCmd.MarkFlagRequired("model-b"); err != nil {
		fmt.Printf("error setting flag required %v", err)
		os.Exit(1)
	}
}
