package list

import (
	"github.com/jmoiron/sqlx"
	"tyre-match-backend/pkg/go-migrations/builder"
)

type CreateTyreModelsTable struct{}

func (m *CreateTyreModelsTable) GetName() string { return "CreateTyreModelsTable" }

func (m *CreateTyreModelsTable) Up(con *sqlx.DB) {
	table := builder.NewTable("tyre_models", con)
	table.Column("id").Type("INT UNSIGNED").NotNull().Autoincrement()
	table.PrimaryKey("id")
	table.String("manufacturer", 255).NotNull()
	table.String("model_name", 255).NotNull()
	table.String("category", 255).Nullable()
	table.String("vehicle_type", 255).Nullable()
	table.Integer("width_mm").Nullable()
	table.Integer("aspect_ratio").Nullable()
	table.Integer("rim_diameter_inches").Nullable()
	table.Integer("groove_count").Nullable()
	table.String("pattern_type", 255).Nullable()
	table.WithTimestamps()
	table.MustExec()
}

func (m *CreateTyreModelsTable) Down(con *sqlx.DB) {
	builder.DropTable("tyre_models", con).MustExec()

}
