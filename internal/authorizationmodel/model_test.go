package authorizationmodel_test

import (
	"testing"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/openfga/pkg/typesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/openfga/cli/internal/authorizationmodel"
)

const (
	modelID        = "01GVKXGDCV2SMG6TRE9NMBQ2VG"
	typeName       = "user"
	modelCreatedAt = "2023-03-16 00:35:51.835 +0000 UTC"
)

func TestReadingInvalidModelFromInvalidJSON(t *testing.T) {
	t.Parallel()

	modelString := "{bad_json"

	model := authorizationmodel.AuthzModel{}

	err := model.ReadFromJSONString(modelString)
	if err == nil {
		t.Errorf("Expected error, got none")
	}
}

func TestReadingValidModelFromJSON(t *testing.T) {
	t.Parallel()

	modelString := `{"id":"01GVKXGDCV2SMG6TRE9NMBQ2VG","schema_version":"1.1","type_definitions":[{"type":"user"}]}`

	model := authorizationmodel.AuthzModel{}

	err := model.ReadFromJSONString(modelString)
	if err != nil {
		t.Errorf("Got error when reading a valid model %v", err)
	}

	if model.GetSchemaVersion() != typesystem.SchemaVersion1_1 {
		t.Errorf("Expected %v to equal %v", model.GetSchemaVersion(), typesystem.SchemaVersion1_1)
	}

	if model.GetID() != modelID {
		t.Errorf("Expected %v to equal %v", model.GetID(), modelID)
	}

	if model.CreatedAt.String() != modelCreatedAt {
		t.Errorf("Expected %v to equal %v", model.GetCreatedAt().String(), modelCreatedAt)
	}

	if model.GetTypeDefinitions()[0].GetType() != typeName {
		t.Errorf("Expected %v to equal %v", model.GetTypeDefinitions()[0].GetType(), typeName)
	}
}

func TestReadingValidModelFromDSL(t *testing.T) {
	t.Parallel()

	model := authorizationmodel.AuthzModel{}
	if err := model.ReadFromDSLString(`model
  schema 1.1

type user
`); err != nil {
		t.Errorf("Got error when parsing a valid model %v", err)
	}

	if model.GetSchemaVersion() != typesystem.SchemaVersion1_1 {
		t.Errorf("Expected %v to equal %v", model.GetSchemaVersion(), typesystem.SchemaVersion1_1)
	}

	if model.GetTypeDefinitions()[0].GetType() != typeName {
		t.Errorf("Expected %v to equal %v", model.GetTypeDefinitions()[0].GetType(), typeName)
	}
}

func TestDisplayAsJsonWithFields(t *testing.T) {
	t.Parallel()

	typeDefs := []openfga.TypeDefinition{{
		Type: typeName,
	}}
	model := authorizationmodel.AuthzModel{
		SchemaVersion:   openfga.PtrString(typesystem.SchemaVersion1_1),
		ID:              openfga.PtrString(modelID),
		TypeDefinitions: &typeDefs,
	}

	jsonModel1 := model.DisplayAsJSON([]string{"model", "id", "created_at"})
	if jsonModel1.GetSchemaVersion() != typesystem.SchemaVersion1_1 {
		t.Errorf("Expected %v to equal %v", jsonModel1.GetSchemaVersion(), typesystem.SchemaVersion1_1)
	}

	if jsonModel1.GetID() != modelID {
		t.Errorf("Expected %v to equal %v", jsonModel1.GetID(), modelID)
	}

	if jsonModel1.GetCreatedAt().String() != modelCreatedAt {
		t.Errorf("Expected %v to equal %v", jsonModel1.GetCreatedAt().String(), modelCreatedAt)
	}

	if jsonModel1.GetTypeDefinitions()[0].GetType() != typeName {
		t.Errorf("Expected %v to equal %v", jsonModel1.GetTypeDefinitions()[0].GetType(), typeName)
	}

	jsonModel2 := model.DisplayAsJSON([]string{"id", "created_at"})
	if jsonModel2.GetSchemaVersion() != "" {
		t.Errorf("Expected %v to be empty", jsonModel2.GetSchemaVersion())
	}

	if jsonModel2.GetID() != modelID {
		t.Errorf("Expected %v to equal %v", jsonModel2.GetID(), modelID)
	}

	if jsonModel1.GetCreatedAt().String() != modelCreatedAt {
		t.Errorf("Expected %v to equal %v", jsonModel2.GetCreatedAt().String(), modelCreatedAt)
	}

	if jsonModel2.GetTypeDefinitions() != nil {
		t.Errorf("Expected %v to equal nil", jsonModel2.GetTypeDefinitions())
	}
}

func TestDiffDSL(t *testing.T) {
	testcases := []struct {
		_name  string
		modelA string
		modelB string
		diff   string
	}{
		{
			_name: "should provide no diff if equal",
			modelA: `model
schema 1.1

type user`,
			modelB: `model
schema 1.1

type user`,
			diff: "",
		},
		{
			_name: "should provide no diff if different order but equal",
			modelA: `model
schema 1.1

type user
type document
  relations 
    define owner: [user]`,
			modelB: `model
schema 1.1

type user
type document
  relations 
    define owner: [user]`,
			diff: "",
		},
		{
			_name: "should provide diff when different",
			modelA: `model
schema 1.1

type user`,
			modelB: `model
schema 1.1

type user
type document
  relations
    define owner: [user]`,
			diff: `model
  schema 1.1

type user

\033[32m+type document
  relations
    define owner: [user]

\033[0m`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc._name, func(t *testing.T) {
			modelA := &authorizationmodel.AuthzModel{}

			err := modelA.ReadFromDSLString(tc.modelA)
			require.NoError(t, err)

			modelB := &authorizationmodel.AuthzModel{}

			err = modelB.ReadFromDSLString(tc.modelB)
			require.NoError(t, err)

			diff, err := modelA.DiffDSL(modelB)

			require.NoError(t, err)

			assert.Equal(t, tc.diff, diff)
		})
	}
}
