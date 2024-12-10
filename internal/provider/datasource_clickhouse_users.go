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
	_ datasource.DataSource              = &clickhouseUsersDataSource{}
	_ datasource.DataSourceWithConfigure = &clickhouseUsersDataSource{}
)

// clickhouseUsersDataSource is the data source implementation.
type clickhouseUsersDataSource struct {
	client clickhouse.Conn
}

// clickhouseUsersDataSourceModel maps the data source schema data.
type clickhouseUsersDataSourceModel struct {
	Users []types.String `tfsdk:"users"`
}

// Metadata returns the data source type name.
func (d *clickhouseUsersDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "clickhouse_users"
}

// Schema defines the schema for the data source.
func (d *clickhouseUsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"users": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of Users in the ClickHouse server.",
			},
		},
	}
}

// Read performs the read operation for the data source.
func (d *clickhouseUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clickhouseUsersDataSourceModel

	rows, err := d.client.Query(ctx, "SHOW USERS")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list Users",
			"An error occurred while listing the Users: "+err.Error(),
		)
		return
	}
	defer rows.Close()

	var Users []types.String
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			resp.Diagnostics.AddError(
				"Unable to read database name",
				"An error occurred while reading the database name: "+err.Error(),
			)
			return
		}
		Users = append(Users, types.StringValue(database))
	}

	state.Users = Users

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure configures the data source with the provider data.
func (d *clickhouseUsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
