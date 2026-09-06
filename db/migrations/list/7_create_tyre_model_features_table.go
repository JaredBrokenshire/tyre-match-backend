package list

import (
	"github.com/jmoiron/sqlx"
	"tyre-match-backend/pkg/go-migrations/builder"
)

type CreateTyreModelFeaturesTable struct{}

func (m *CreateTyreModelFeaturesTable) GetName() string { return "CreateTyreModelFeaturesTable" }

func (m *CreateTyreModelFeaturesTable) Up(con *sqlx.DB) {
	table := builder.NewTable("tyre_model_features", con)
	table.Column("id").Type("INT UNSIGNED").NotNull().Autoincrement()
	table.PrimaryKey("id")
	table.Column("tyre_model_id").Type("INT UNSIGNED").NotNull().Unique()
	table.Column("feature_json").Type("LONGTEXT").NotNull()
	table.WithTimestamps()
	table.MustExec()
}

func (m *CreateTyreModelFeaturesTable) Down(con *sqlx.DB) {
	builder.DropTable("tyre_model_features", con).MustExec()
}
