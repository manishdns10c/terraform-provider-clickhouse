package provider

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clickhouseDatabasesDataSource{}
	_ datasource.DataSourceWithConfigure = &clickhouseDatabasesDataSource{}
)

// clickhouseDatabasesDataSource is the data source implementation.
type clickhouseDatabasesDataSource struct {
	client clickhouse.Conn
}

// clickhouseDatabasesDataSourceModel maps the data source schema data.
type clickhouseDatabasesDataSourceModel struct {
	Databases []types.String `tfsdk:"databases"`
}

// Metadata returns the data source type name.
func (d *clickhouseDatabasesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "clickhouse_databases"
}

// Schema defines the schema for the data source.
func (d *clickhouseDatabasesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"databases": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of databases in the ClickHouse server.",
			},
		},
	}
}

// Read performs the read operation for the data source.
func (d *clickhouseDatabasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clickhouseDatabasesDataSourceModel

	rows, err := d.client.Query(ctx, "SHOW DATABASES")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list databases",
			"An error occurred while listing the databases: "+err.Error(),
		)
		return
	}
	defer rows.Close()

	var databases []types.String
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			resp.Diagnostics.AddError(
				"Unable to read database name",
				"An error occurred while reading the database name: "+err.Error(),
			)
			return
		}
		databases = append(databases, types.StringValue(database))
	}

	state.Databases = databases

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure configures the data source with the provider data.
func (d *clickhouseDatabasesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(clickhouse.Conn)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected clickhouse.Conn, got something else",
		)
		return
	}

	d.client = client
}
