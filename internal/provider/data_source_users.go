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
	"github.com/keeper-security/keeper-sdk-golang/api"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
	"reflect"
	"strings"
	"terraform-provider-keeper/internal/model"
)

var (
	_ datasource.DataSource = &usersDataSource{}
)

func newUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSourceModel struct {
	IsActive     types.Bool            `tfsdk:"is_active"`
	Filter       *model.FilterCriteria `tfsdk:"filter"`
	Emails       types.Set             `tfsdk:"emails"`
	NodeCriteria *model.NodeCriteria   `tfsdk:"nodes"`
	Users        []*model.UserModel    `tfsdk:"users"`
}

type usersDataSource struct {
	users enterprise.IEnterpriseEntity[enterprise.IUser, int64]
	nodes enterprise.IEnterpriseEntity[enterprise.INode, int64]
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
			"nodes": schema.SingleNestedAttribute{
				Attributes:  model.NodeCriteriaAttributes,
				Optional:    true,
				Description: "Search By node filter",
			},
			"filter": schema.SingleNestedAttribute{
				Attributes:  model.FilterCriteriaAttributes,
				Optional:    true,
				Description: "Search By field filter",
			},
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: model.UserSchemaAttributes,
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
		d.nodes = ed.Nodes()
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
		var activeFlag = enterprise.UserStatus_Inactive
		if uq.IsActive.ValueBool() {
			activeFlag = enterprise.UserStatus_Active
		}
		activeMatcher = func(user enterprise.IUser) bool {
			return user.Status() == activeFlag
		}
	}

	var userMatcher model.Matcher
	var nodeMatcher model.NodeMatcher
	var diags diag.Diagnostics
	if uq.Filter != nil {
		userMatcher, diags = model.GetFieldMatcher(uq.Filter, reflect.TypeOf((*model.UserModel)(nil)))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else if uq.NodeCriteria != nil {
		nodeMatcher, diags = model.GetNodeMatcher(uq.NodeCriteria, d.nodes)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else if !uq.Emails.IsNull() {
		var us = api.NewSet[string]()

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
				us.Add(strings.ToLower(strv.ValueString()))
			}
		}

		userMatcher = func(mm interface{}) (b bool) {
			var ok bool
			var m1 *model.UserModel
			if m1, ok = mm.(*model.UserModel); ok {
				if !m1.Username.IsNull() {
					var username = strings.ToLower(m1.Username.ValueString())
					b = us.Has(username)
				}
			}
			return
		}
	}

	state.Users = make([]*model.UserModel, 0)
	d.users.GetAllEntities(func(u enterprise.IUser) (resume bool) {
		resume = true
		if !activeMatcher(u) {
			return
		}
		if nodeMatcher != nil {
			if !nodeMatcher(u.NodeId()) {
				return
			}
		}
		var um = new(model.UserModel)
		um.FromKeeper(u)
		if userMatcher != nil {
			if !userMatcher(um) {
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
