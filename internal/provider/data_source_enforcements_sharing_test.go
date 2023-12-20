package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceEnforcementsSharing(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnforcementsSharingDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_enforcements_sharing.test", "restrict_import", "true"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_sharing.test", "restrict_export", "true"),
				),
			},
		}})
}

const testAccEnforcementsSharingDataSourceConfig = `
data "keeper_enforcements_sharing" "test" {
  restrict_import = true
  restrict_export = true
}
`
