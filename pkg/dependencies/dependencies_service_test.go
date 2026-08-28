package dependencies_test

import (
	"github.com/joho/godotenv"
	"testing"
	"tyre-match-backend/pkg/dependencies"
	"tyre-match-backend/tests/helpers"
)

func TestNewDependencyService(t *testing.T) {

	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file")
	}

	ds := dependencies.NewDependencyService(helpers.MockDb())

	ds.GetDB()

}
