package archive

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testProviders = map[string]terraform.ResourceProvider{
	"archive": Provider(),
}

func TestMain(m *testing.M) {
	acctest.UseBinaryDriver("archive", Provider)
	resource.TestMain(m)
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
