package archive

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns the resources and schema defined for this provider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"archive_file": dataSourceFile(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"archive_file": schema.DataSourceResourceShim(
				"archive_file",
				dataSourceFile(),
			),
		},
	}
}
