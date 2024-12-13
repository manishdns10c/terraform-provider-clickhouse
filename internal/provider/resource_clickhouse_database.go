package provider

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &clickhouseDatabaseResource{}
	_ resource.ResourceWithConfigure = &clickhouseDatabaseResource{}
)

// clickhousedatabaseResource is the resource implementation.
type clickhouseDatabaseResource struct {
	client clickhouse.Conn
}

// clickhousedatabaseResourceModel maps the resource schema data.
type clickhouseDatabaseResourceModel struct {
	Database types.String `tfsdk:"database"`
}

// Metadata returns the resource type name.
func (r *clickhouseDatabaseResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "clickhouse_database"
}

// Schema defines the schema for the resource.
func (r *clickhouseDatabaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"database": schema.StringAttribute{
				Required:    true,
				Description: "The name of the ClickHouse database.",
			},
		},
	}
}

// Create handles the creation of the resource.
func (r *clickhouseDatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan clickhouseDatabaseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createDatabseQuery := fmt.Sprintf(
		"CREATE DATABASE %s",
		plan.Database.ValueString(),
	)

	if err := r.client.Exec(ctx, createDatabseQuery); err != nil {
		resp.Diagnostics.AddError(
			"Error creating ClickHouse database",
			"Could not create ClickHouse database, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource data.
func (r *clickhouseDatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state clickhouseDatabaseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Properly format the query to check if the database exists, enclosing the databasename in single quotes
	query := fmt.Sprintf("SELECT count() > 0 FROM system.databases WHERE name='%s'", state.Database.ValueString())
	var exists bool
	err := r.client.QueryRow(ctx, query).Scan(&exists)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ClickHouse database",
			"Could not read ClickHouse database, unexpected error: "+err.Error(),
		)
		return
	}

	if !exists {
		resp.Diagnostics.AddError(
			"Database does not exist",
			"The ClickHouse database "+state.Database.ValueString()+" does not exist.",
		)
		return
	}

	// If the database exists, set the state and continue
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
// Update handles updating the resource for a database.
// In this case, we simply return an error indicating that updates are not permitted.
func (r *clickhouseDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Create a new diagnostics entry indicating that updates are not permitted
	resp.Diagnostics.AddError(
		"Update Not Permitted",
		"Updating an existing ClickHouse database is not permitted. Please recreate the database instead of attempting to update it.",
	)

	// No further logic is needed as we're just returning an error message.
}

// Delete handles deleting the resource.
func (r *clickhouseDatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clickhouseDatabaseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteDatabaseQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", state.Database.ValueString())

	if err := r.client.Exec(ctx, deleteDatabaseQuery); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ClickHouse database",
			"Could not delete ClickHouse database, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure configures the resource with the provider data.
func (r *clickhouseDatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(clickhouse.Conn)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected clickhouse.Conn, got something else",
		)
		return
	}

	r.client = client
}
