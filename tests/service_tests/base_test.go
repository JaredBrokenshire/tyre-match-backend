package service_tests

import (
	"github.com/jinzhu/gorm"
	"os"
	"testing"
	"tyre-match-backend/tests/helpers"
)

var (
	ts *helpers.TestServer
)

func TestMain(m *testing.M) {
	ts = helpers.NewTestServer("../../.env", false)
	// Close the database connection
	defer ts.S.Db.Close()

	// Run the test
	code := m.Run()

	os.Exit(code)
}

// / CheckDatabase runs a query and returns a count of rows found
func CheckDatabase(dbQuery *helpers.DatabaseCheck) int {
	var resultCount int
	var queryScopes []func(db *gorm.DB) *gorm.DB

	modelScope := func(db *gorm.DB) *gorm.DB { return db.Where(dbQuery.Model) }
	queryScopes = append(queryScopes, modelScope)

	if dbQuery.Scope != nil {
		queryScopes = append(queryScopes, dbQuery.Scope)
	}

	db := ts.S.Db
	// Enable debugging on this query if set in the dbQuery
	if dbQuery.DebugQuery {
		db = ts.S.Db.Debug()
	}

	db.Model(dbQuery.Model).Scopes(queryScopes...).Count(&resultCount)
	return resultCount
}
