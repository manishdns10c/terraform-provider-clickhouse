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
	_ resource.Resource              = &clickhouseUserResource{}
	_ resource.ResourceWithConfigure = &clickhouseUserResource{}
)

// clickhouseUserResource is the resource implementation.
type clickhouseUserResource struct {
	client clickhouse.Conn
}

// clickhouseUserResourceModel maps the resource schema data.
type clickhouseUserResourceModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the resource type name.
func (r *clickhouseUserResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "clickhouse_user"
}

// Schema defines the schema for the resource.
func (r *clickhouseUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Required:    true,
				Description: "The name of the ClickHouse user.",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Description: "The password of the ClickHouse user.",
				Sensitive:   true,
			},
		},
	}
}

// Create handles the creation of the resource.
func (r *clickhouseUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan clickhouseUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createUserQuery := fmt.Sprintf(
		"CREATE USER %s IDENTIFIED BY '%s'",
		plan.Username.ValueString(),
		plan.Password.ValueString(),
	)

	if err := r.client.Exec(ctx, createUserQuery); err != nil {
		resp.Diagnostics.AddError(
			"Error creating ClickHouse user",
			"Could not create ClickHouse user, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource data.
func (r *clickhouseUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state clickhouseUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Properly format the query to check if the user exists, enclosing the username in single quotes
	query := fmt.Sprintf("SELECT count() > 0 FROM system.users WHERE name = '%s'", state.Username.ValueString())
	var exists bool
	err := r.client.QueryRow(ctx, query).Scan(&exists)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ClickHouse user",
			"Could not read ClickHouse user, unexpected error: "+err.Error(),
		)
		return
	}

	if !exists {
		resp.Diagnostics.AddError(
			"User does not exist",
			"The ClickHouse user "+state.Username.ValueString()+" does not exist.",
		)
		return
	}

	// If the user exists, set the state and continue
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *clickhouseUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan clickhouseUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateUserQuery := fmt.Sprintf(
		"ALTER USER %s IDENTIFIED BY '%s'",
		plan.Username.ValueString(),
		plan.Password.ValueString(),
	)

	if err := r.client.Exec(ctx, updateUserQuery); err != nil {
		resp.Diagnostics.AddError(
			"Error updating ClickHouse user",
			"Could not update ClickHouse user, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *clickhouseUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clickhouseUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteUserQuery := fmt.Sprintf("DROP USER %s", state.Username.ValueString())

	if err := r.client.Exec(ctx, deleteUserQuery); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ClickHouse user",
			"Could not delete ClickHouse user, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure configures the resource with the provider data.
func (r *clickhouseUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
