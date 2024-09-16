package authorizationmodel

import (
	"fmt"
	"strings"

	"github.com/openfga/cli/internal/cmdutils"
	openfga "github.com/openfga/go-sdk"
	"github.com/spf13/cobra"
)

func isFileBased(model string) bool {
	return strings.HasSuffix(model, ".fga") ||
		strings.HasSuffix(model, ".fga.yaml") ||
		strings.HasSuffix(model, ".mod")
}

func Read(cmd *cobra.Command, flagName string) (*AuthzModel, error) {
	modelFlagValue, err := cmd.Flags().GetString(flagName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model name due to %w", err)
	}

	authModel := &AuthzModel{}

	if isFileBased(modelFlagValue) {
		var model string

		format := ModelFormatDefault

		if err := ReadFromFile(
			modelFlagValue,
			&model,
			&format,
			openfga.PtrString(""),
		); err != nil {
			return nil, err //nolint:wrapcheck
		}

		if err := authModel.ReadModelFromString(model, format); err != nil {
			return nil, err //nolint:wrapcheck
		}
	} else {
		clientConfig := cmdutils.GetClientConfig(cmd)

		fgaClient, err := clientConfig.GetFgaClient()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize FGA Client due to %w", err)
		}

		if modelFlagValue != "latest" {
			clientConfig.AuthorizationModelID = modelFlagValue
		}

		response, err := ReadFromStore(clientConfig, fgaClient)
		if err != nil {
			return nil, err //nolint:wrapcheck
		}

		authModel.Set(*response.AuthorizationModel)
	}

	return authModel, nil
}
