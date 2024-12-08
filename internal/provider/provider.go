package provider

import (
	"context"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &clickhouseProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &clickhouseProvider{
			version: version,
		}
	}
}

// clickhouseProvider is the provider implementation.
type clickhouseProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// clickhouseProviderModel maps provider schema data to a Go type.
type clickhouseProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *clickhouseProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "clickhouse"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
// Schema defines the provider-level schema for configuration data.
func (p *clickhouseProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required:    true,
				Description: "The hostname or IP address of the ClickHouse server.",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "The username for accessing the ClickHouse server.",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The password for accessing the ClickHouse server.",
			},
		},
	}
}

// Configure prepares a clickhouse API client for data sources and resources.
func (p *clickhouseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config clickhouseProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("CLICKHOUSE_HOST")
	username := os.Getenv("CLICKHOUSE_USERNAME")
	password := os.Getenv("CLICKHOUSE_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing ClickHouse Host",
			"The provider cannot create the ClickHouse client as there is a missing or empty value for the ClickHouse host. "+
				"Set the host value in the configuration or use the CLICKHOUSE_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing ClickHouse Username",
			"The provider cannot create the ClickHouse client as there is a missing or empty value for the ClickHouse username. "+
				"Set the username value in the configuration or use the CLICKHOUSE_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing ClickHouse Password",
			"The provider cannot create the ClickHouse client as there is a missing or empty value for the ClickHouse password. "+
				"Set the password value in the configuration or use the CLICKHOUSE_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{host},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create ClickHouse Client",
			"An unexpected error occurred when creating the ClickHouse client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"ClickHouse Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
// DataSources defines the data sources implemented in the provider.
func (p *clickhouseProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &clickhouseDatabasesDataSource{}
		},
	}
}

// Resources defines the resources implemented in the provider.
func (p *clickhouseProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
