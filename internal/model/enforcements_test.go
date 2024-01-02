package model

import (
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/keeper-security/keeper-sdk-golang/api"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestEnforcements(t *testing.T) {
	var descriptionEnforcements = api.NewSet[string]()
	for k := range enforcementDescriptions {
		descriptionEnforcements.Add(k)
	}

	ea, diags := EnforcementsDataSourceAttributes()
	assert.Assert(t, len(diags) == 0)
	var modelEnforcements = api.NewSet[string]()
	for k, v := range ea {
		var e, ok = v.(dsschema.SingleNestedAttribute)
		assert.Assert(t, ok)
		descriptionEnforcements.Delete(k)
		for k = range e.Attributes {
			modelEnforcements.Add(k)
		}
	}

	var sdkEnforcements = api.NewSet[string]()
	enterprise.AvailableRoleEnforcements(func(x enterprise.IEnforcement) bool {
		sdkEnforcements.Add(strings.ToLower(x.Name()))
		return true
	})

	var excludedEnforcements = []string{"two_factor_by_ip"}
	sdkEnforcements.Difference(excludedEnforcements)
	descriptionEnforcements.Difference(excludedEnforcements)

	var result = api.MakeSet[string](modelEnforcements.ToArray())
	result.Difference(sdkEnforcements.ToArray())
	assert.Assert(t, len(result) == 0, "Model defines unsupported enforcement(s): %s", strings.Join(result.ToArray(), ", "))

	result = api.MakeSet[string](sdkEnforcements.ToArray())
	result.Difference(modelEnforcements.ToArray())
	assert.Assert(t, len(result) == 0, "Model misses the following SDK enforcement(s): %s", strings.Join(result.ToArray(), ", "))

	result = api.MakeSet[string](modelEnforcements.ToArray())
	result.Difference(descriptionEnforcements.ToArray())
	assert.Assert(t, len(result) == 0, "Model misses the following enforcement description(s): %s", strings.Join(result.ToArray(), ", "))

	result = api.MakeSet[string](descriptionEnforcements.ToArray())
	result.Difference(modelEnforcements.ToArray())
	assert.Assert(t, len(result) == 0, "The following enforcement description(s) are unused: %s", strings.Join(result.ToArray(), ", "))
}
