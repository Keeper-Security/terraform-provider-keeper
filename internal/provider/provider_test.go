package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
)

var management enterprise.IEnterpriseManagement
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"keeper": providerserver.NewProtocol6WithError(New("test", &management)()),
}
