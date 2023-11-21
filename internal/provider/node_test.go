package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"gotest.tools/assert"
	"testing"
)

func TestDataSource(t *testing.T) {
	var resp = new(datasource.SchemaResponse)
	var ds = new(teamDataSource)
	ds.Schema(mockContext{}, datasource.SchemaRequest{}, resp)
	assert.Assert(t, resp.Schema.Description)
}
