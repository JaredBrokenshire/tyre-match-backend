package list

import (
	"github.com/jmoiron/sqlx"
	"tyre-match-backend/pkg/go-migrations/builder"
)

type CreateTyreImpressionsTable struct{}

func (m *CreateTyreImpressionsTable) GetName() string {
	return "CreateTyreImpressionsTable"
}

func (m *CreateTyreImpressionsTable) Up(con *sqlx.DB) {
	table := builder.NewTable("tyre_impressions", con)
	table.Column("id").Type("INT UNSIGNED").NotNull().Autoincrement()
	table.PrimaryKey("id")
	table.Column("status").Type("ENUM('Uploaded','Processing','Processed','Matched','Failed')").NotNull().Default("Uploaded")
	table.Column("pixels_per_inch").Type("FLOAT").NotNull()
	table.Column("edge_density").Type("FLOAT").Nullable()
	table.Column("void_ratio").Type("FLOAT").Nullable()
	table.Integer("groove_count").Nullable()
	table.WithTimestamps()
	table.MustExec()
}

func (m *CreateTyreImpressionsTable) Down(con *sqlx.DB) {
	builder.DropTable("tyre_impressions", con).MustExec()
}
