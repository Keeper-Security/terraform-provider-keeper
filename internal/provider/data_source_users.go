package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"reflect"
	"strings"
)

var (
	_ datasource.DataSource = &usersDataSource{}
)

type usersDataSourceModel struct {
	IsActive types.Bool      `tfsdk:"is_active"`
	Filter   *filterCriteria `tfsdk:"filter"`
	Emails   types.Set       `tfsdk:"emails"`
	Nodes    types.Set       `tfsdk:"nodes"`
	Users    []*userModel    `tfsdk:"users"`
}

type usersDataSource struct {
	users enterprise.IEnterpriseEntity[enterprise.IUser, int64]
}

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}
func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"is_active": schema.BoolAttribute{
				Optional:    true,
				Description: "Filter Active or Pending users",
			},
			"emails": schema.SetAttribute{
				Optional:    true,
				Description: "List of emails",
				ElementType: types.StringType,
			},
			"nodes": schema.SetAttribute{
				Optional:    true,
				Description: "List of node IDs",
				ElementType: types.Int64Type,
			},
			"filter": schema.SingleNestedAttribute{
				Attributes:  filterCriteriaAttributes,
				Optional:    true,
				Description: "Search By field filter",
			},
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: userSchemaAttributes,
				},
			},
		},
	}
}

func (d *usersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if ed, ok := req.ProviderData.(enterprise.IEnterpriseData); ok {
		d.users = ed.Users()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected \"IEnterpriseLoader\", got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var uq usersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &uq)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = uq

	var activeMatcher func(user enterprise.IUser) bool
	if uq.IsActive.IsNull() {
		activeMatcher = func(_ enterprise.IUser) bool {
			return true
		}
	} else {
		var activeFlag = "inactive"
		if uq.IsActive.ValueBool() {
			activeFlag = "active"
		}
		activeMatcher = func(user enterprise.IUser) bool {
			return user.Status() == activeFlag
		}
	}

	var cb matcher
	var diags diag.Diagnostics
	if uq.Filter != nil {
		cb, diags = getFieldMatcher(uq.Filter, reflect.TypeOf((*userModel)(nil)))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else if !uq.Nodes.IsNull() {
		var nodes = make(map[int64]bool)
		var sv types.Set
		sv, diags = uq.Nodes.ToSetValue(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range sv.Elements() {
			var ok bool
			var i64v types.Int64
			if i64v, ok = v.(types.Int64); ok {
				nodes[i64v.ValueInt64()] = true
			}
		}
		cb = func(mm interface{}) (b bool) {
			var ok bool
			var m1 *userModel
			if m1, ok = mm.(*userModel); ok {
				if !m1.NodeId.IsNull() {
					_, b = nodes[m1.NodeId.ValueInt64()]
				}
			}
			return
		}
	} else if !uq.Emails.IsNull() {
		var nodes = make(map[string]bool)
		var sv types.Set
		sv, diags = uq.Emails.ToSetValue(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range sv.Elements() {
			var ok bool
			var strv types.String
			if strv, ok = v.(types.String); ok {
				nodes[strings.ToLower(strv.ValueString())] = true
			}
		}
		cb = func(mm interface{}) (b bool) {
			var ok bool
			var m1 *userModel
			if m1, ok = mm.(*userModel); ok {
				if !m1.Username.IsNull() {
					_, b = nodes[strings.ToLower(m1.Username.ValueString())]
				}
			}
			return
		}
	}

	d.users.GetAllEntities(func(u enterprise.IUser) (resume bool) {
		resume = true
		if !activeMatcher(u) {
			return
		}
		var um = new(userModel)
		um.fromKeeper(u)
		if cb != nil {
			if !cb(um) {
				return
			}
		}
		state.Users = append(state.Users, um)

		return
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *usersDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("emails"),
			path.MatchRoot("nodes"),
			path.MatchRoot("filter"),
		),
	}
}
