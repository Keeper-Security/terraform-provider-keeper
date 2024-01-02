package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccDataSourceEnforcementsLogin(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccEnforcementsLoginErrorDataSourceConfig,
				ExpectError: regexp.MustCompile(`Boolean enforcement`),
			},
			{
				Config: testAccEnforcementsLoginDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_enforcements_login.test", "restrict_windows_fingerprint", "true"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_login.test", "master_password_maximum_days_before_change", "20"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_login.test", "master_password_minimum_length", "20"),
				),
			},
		}})
}

const testAccEnforcementsLoginDataSourceConfig = `
data "keeper_enforcements_login" "test" {
  restrict_windows_fingerprint = true
  master_password_maximum_days_before_change = 20
  master_password_minimum_length = 20
}
`

const testAccEnforcementsLoginErrorDataSourceConfig = `
data "keeper_enforcements_login" "test" {
  restrict_windows_fingerprint = false
}
`
