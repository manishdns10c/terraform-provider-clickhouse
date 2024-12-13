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
	_ datasource.DataSource              = &clickhouseRolesDataSource{}
	_ datasource.DataSourceWithConfigure = &clickhouseRolesDataSource{}
)

// clickhouseUsersDataSource is the data source implementation.
type clickhouseRolesDataSource struct {
	client clickhouse.Conn
}

// clickhouseUsersDataSourceModel maps the data source schema data.
type clickhouseRolesDataSourceModel struct {
	Roles []types.String `tfsdk:"roles"`
}

// Metadata returns the data source type name.
func (d *clickhouseRolesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "clickhouse_roles"
}

// Schema defines the schema for the data source.
func (d *clickhouseRolesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"roles": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of Roles in the ClickHouse server.",
			},
		},
	}
}

// Read performs the read operation for the data source.
func (d *clickhouseRolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clickhouseRolesDataSourceModel

	rows, err := d.client.Query(ctx, "SHOW ROLES")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list Roles",
			"An error occurred while listing the roles: "+err.Error(),
		)
		return
	}
	defer rows.Close()

	var Roles []types.String
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			resp.Diagnostics.AddError(
				"Unable to read database name",
				"An error occurred while reading the database name: "+err.Error(),
			)
			return
		}
		Roles = append(Roles, types.StringValue(database))
	}

	state.Roles = Roles

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure configures the data source with the provider data.
func (d *clickhouseRolesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
