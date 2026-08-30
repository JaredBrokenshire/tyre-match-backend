package process

import (
	"tyre-match-backend/db/migrations/list"
	"tyre-match-backend/pkg/go-migrations"
	"tyre-match-backend/pkg/go-migrations/store"
)

func Run() {
	go_migrations.Run(getMigrationsList())
}

func getMigrationsList() []store.Migratable {
	return []store.Migratable{
		&list.CreateTyreModelsTable{},
		&list.CreateTyreImpressionsTable{},
		&list.CreateFilesTable{},
		&list.UpdateFilesFileTypeAddBinary{},
		&list.UpdateFilesFileTypeAddNormalised{},
		&list.UpdateTyreImpressionsAddROI{},
	}
}
