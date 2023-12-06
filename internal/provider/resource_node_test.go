package provider

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccNodeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccCreateSsoNodeConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kepr_node.SSO", "name", "SSO"),
					resource.TestCheckResourceAttr("kepr_node.SSO", "restrict_visibility", "true"),
				),
			},
			// Import
			{
				Config:        testAccImportSubnode(),
				ImportState:   true,
				ResourceName:  "kepr_node.subnode",
				ImportStateId: "5299989643274",
				ImportStateCheck: func(states []*terraform.InstanceState) (err error) {
					var state = states[0]
					var eId int
					var value string
					var ok bool
					if value, ok = state.Attributes["parent_id"]; ok {
						if eId, err = strconv.Atoi(value); err != nil {
							return
						}
					} else {
						err = errors.New("\"parent_id\" attribute has not been set")
						return
					}
					if eId != 5299989643266 {
						err = errors.New(fmt.Sprintf("\"parent_id\" value is incorrect. Expected: 5299989643266; Got: %d", eId))
						return
					}
					if value, ok = state.Attributes["name"]; ok {
						if value != "Subnode" {
							err = errors.New(fmt.Sprintf("\"name\" value is incorrect. Expected: \"Subnode\"; Got: %s", value))
							return
						}
					} else {
						err = errors.New("\"name\" attribute has not been set")
						return
					}

					return
				},
				ImportStatePersist: true,
			},
			// Update and Read testing
			{
				Config: testAccImportSubnode(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kepr_node.subnode", "name", "Subnodes"),
				),
			},
			//// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCreateSsoNodeConfig() string {
	return `
data "kepr_node" "root" {
	is_root = true
}
resource "kepr_node" "SSO" {
	name = "SSO"
	parent_id = data.kepr_node.root.node_id
	restrict_visibility = true
}
`
}

func testAccImportSubnode() string {
	return `
resource "kepr_node" "subnode" {
	name = "Subnodes"
}
`
}
