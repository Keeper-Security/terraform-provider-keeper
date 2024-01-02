package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccDataSourceEnforcementsRecordTypes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccEnforcementsErrorRecordTypesDataSourceConfig,
				ExpectError: regexp.MustCompile("Record type name are case-sensitive"),
			},
			{
				Config: testAccEnforcementsRecordTypesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("data.keeper_enforcements_record_types.test", "restrict_record_types.*", "contact"),
				),
			},
		}})
}

const testAccEnforcementsRecordTypesDataSourceConfig = `
data "keeper_enforcements_record_types" "test" {
	restrict_record_types = ["contact", "bankAccount"]
}
`
const testAccEnforcementsErrorRecordTypesDataSourceConfig = `
data "keeper_enforcements_record_types" "test" {
	restrict_record_types = ["contact", "bankaccount"]
}
`
