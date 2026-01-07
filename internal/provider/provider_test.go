// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

//nolint:unparam
func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"archive": providerserver.NewProtocol5WithError(New()),
	}
}
