package archive

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func New() provider.Provider {
	return &archiveProvider{}
}

var _ provider.Provider = (*archiveProvider)(nil)

type archiveProvider struct{}

func (p *archiveProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{}, nil
}

func (p *archiveProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

func (p *archiveProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewArchiveFileResource,
	}
}

func (p *archiveProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewArchiveFileDataSource,
	}
}

func (p *archiveProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "archive"
}
