package store

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slices"
)

type Migratable interface {
	GetName() string
	Up(*sqlx.DB)
	Down(*sqlx.DB)
}

var list []Migratable

func RegisterMigrations(migs []Migratable) {
	list = migs
}

func GetMigrationsList() []Migratable {
	return list
}

func FindMigration(name string) Migratable {
	for _, m := range list {
		if m.GetName() == name {
			return m
		}
	}
	return nil
}

func FilterMigrations(name string) {
	if list[len(list)-1].GetName() == name {
		list = []Migratable{}
	}
	for i, m := range list {
		if m.GetName() == name {
			list = list[i+1:]
		}
	}
}

// FilterToMissingMigrations will loop though the migration in the DB and filter them out of the `list`
func FilterToMissingMigrations(migrationsInDb []string) {
	var newList []Migratable
	// Loop though the list and only add it ot the new list if it DOESN't exist in the list of migrations in DB
	for _, listMigration := range list {
		if slices.Contains(migrationsInDb, listMigration.GetName()) == false {
			newList = append(newList, listMigration)
		}
	}
	// Replace the list
	list = newList
}
