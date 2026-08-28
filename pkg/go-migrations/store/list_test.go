package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockMigratable struct{ name string }

func (m mockMigratable) GetName() string { return m.name }
func (m mockMigratable) Up(*sqlx.DB)     {}
func (m mockMigratable) Down(*sqlx.DB)   {}

func TestFilterToMissingMigrations(t *testing.T) {

	// Fill the current list. Note the incorrect 5 then 4 order. We will presume the DB has 1-4
	list = []Migratable{
		mockMigratable{"migration1"},
		mockMigratable{"migration2"},
		mockMigratable{"migration3"},
		mockMigratable{"migration5"},
		mockMigratable{"migration4"},
		mockMigratable{"migration6"},
	}

	// Pass through migrations 1-4 as what exists in the DB.
	FilterToMissingMigrations([]string{
		"migration1", "migration2", "migration3", "migration4",
	})

	// List should now only have 5 and 6 in.
	assert.Len(t, list, 2, "List should only have two objects in. 5 and 6")
	assert.Equal(t, list[0].GetName(), "migration5")
	assert.Equal(t, list[1].GetName(), "migration6")

}
