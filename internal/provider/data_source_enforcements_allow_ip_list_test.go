package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceEnforcementsAllowIpList(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnforcementsAllowIpListDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_enforcements_allow_ip_list.test", "restrict_ip_addresses.#", "2"),
				),
			},
		}})
}

const testAccEnforcementsAllowIpListDataSourceConfig = `
data "keeper_enforcements_allow_ip_list" "test" {
  restrict_ip_addresses = ["12.13.14.0-12.13.14.255","192.168.1-192.168.1.254"]
}
`
