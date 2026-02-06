// Copyright IBM Corp. 2017, 2026
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func New() provider.Provider {
	return &archiveProvider{}
}

var _ provider.Provider = (*archiveProvider)(nil)

type archiveProvider struct{}

func (p *archiveProvider) Schema(context.Context, provider.SchemaRequest, *provider.SchemaResponse) {
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
