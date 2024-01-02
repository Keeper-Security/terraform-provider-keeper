package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/keeper-security/keeper-sdk-golang/api"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
	"github.com/keeper-security/keeper-sdk-golang/vault"
	"strconv"
	"strings"
	"terraform-provider-keeper/internal/model"
)

func newRoleEnforcementsResource() resource.Resource {
	return &roleEnforcementsResource{}
}

type roleEnforcementsResourceModel struct {
	RoleId       types.Int64                        `tfsdk:"role_id"`
	Enforcements *model.EnforcementsDataSourceModel `tfsdk:"enforcements"`
}

type roleEnforcementsResource struct {
	management enterprise.IEnterpriseManagement
}

func (rmr *roleEnforcementsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_enforcements"
}

func (rmr *roleEnforcementsResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	var diags diag.Diagnostics
	var attributes map[string]schema.Attribute
	attributes, diags = model.EnforcementResourceAttributes()
	resp.Diagnostics.Append(diags...)
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"role_id": schema.Int64Attribute{
				Required:    true,
				Description: "Role ID",
			},
			"enforcements": schema.SingleNestedAttribute{
				Required:    true,
				Attributes:  attributes,
				Description: "Role Enforcements",
			},
		},
	}
}

func (rmr *roleEnforcementsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	mgmt, ok := req.ProviderData.(enterprise.IEnterpriseManagement)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected IEnterpriseManagement, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	rmr.management = mgmt
}

func (rmr *roleEnforcementsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var roleId int
	var err error

	if roleId, err = strconv.Atoi(req.ID); err != nil {
		resp.Diagnostics.AddError(
			"Invalid role ID value.",
			fmt.Sprintf("Interger value expected. Got \"%s\"", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_id"), int64(roleId))...)
}

func (rmr *roleEnforcementsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state roleEnforcementsResourceModel
	var diags = req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var roles = rmr.management.EnterpriseData().Roles()
	var roleId = state.RoleId.ValueInt64()
	var role = roles.GetEntity(roleId)
	if role == nil {
		resp.Diagnostics.AddError(
			"Role not found",
			fmt.Sprintf("Role ID \"%d\" does not exist", roleId),
		)
		return
	}

	state.RoleId = types.Int64Value(roleId)
	var enforcements = make(map[string]string)
	rmr.management.EnterpriseData().RoleEnforcements().GetLinksBySubject(roleId, func(x enterprise.IRoleEnforcement) bool {
		enforcements[strings.ToLower(x.EnforcementType())] = x.Value()
		return true
	})
	state.Enforcements = new(model.EnforcementsDataSourceModel)
	var rts []vault.IRecordType
	if rmr.management.EnterpriseData().RecordTypes() != nil {
		rmr.management.EnterpriseData().RecordTypes().GetAllEntities(func(x vault.IRecordType) bool {
			rts = append(rts, x)
			return true
		})
	}
	state.Enforcements.FromKeeper(enforcements, rts)

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (rmr *roleEnforcementsResource) applyEnforcements(ctx context.Context, enf roleEnforcementsResourceModel) (diags diag.Diagnostics) {
	var roles = rmr.management.EnterpriseData().Roles()
	var roleId = enf.RoleId.ValueInt64()
	var role = roles.GetEntity(roleId)
	if role == nil {
		diags.AddError(
			"Role not found",
			fmt.Sprintf("Role ID \"%d\" does not exist", roleId),
		)
		return
	}

	var origEnforcements = make(map[string]string)
	rmr.management.EnterpriseData().RoleEnforcements().GetLinksBySubject(roleId, func(x enterprise.IRoleEnforcement) bool {
		origEnforcements[strings.ToLower(x.EnforcementType())] = x.Value()
		return true
	})
	tflog.Debug(ctx, fmt.Sprintf("# of orig enforcements: %d", len(origEnforcements)))

	var newEnforcements = make(map[string]string)
	if enf.Enforcements != nil {
		var rts []vault.IRecordType
		if rmr.management.EnterpriseData().RecordTypes() != nil {
			rmr.management.EnterpriseData().RecordTypes().GetAllEntities(func(x vault.IRecordType) bool {
				rts = append(rts, x)
				return true
			})
		}
		enf.Enforcements.ToKeeper(newEnforcements, rts)
	}
	tflog.Debug(ctx, fmt.Sprintf("# of new enforcements: %d", len(newEnforcements)))

	var oldKeys = api.NewSet[string]()
	var newKeys = api.NewSet[string]()
	for k := range origEnforcements {
		oldKeys.Add(k)
	}
	for k := range newEnforcements {
		newKeys.Add(k)
	}
	var enforcements []enterprise.IRoleEnforcement
	var keys = api.NewSet[string]()
	keys.Union(oldKeys.ToArray())
	keys.Difference(newKeys.ToArray())
	if len(keys) > 0 { // remove
		for _, k := range keys.ToArray() {
			tflog.Debug(ctx, fmt.Sprintf("delete enforcements: %s", k))
			enforcements = append(enforcements, enterprise.NewRoleEnforcement(roleId, k))
		}
	}
	keys = api.NewSet[string]()
	keys.Union(newKeys.ToArray())
	keys.Difference(oldKeys.ToArray())

	var ok bool
	var strValue string
	if len(keys) > 0 { // add
		for _, k := range keys.ToArray() {
			if strValue, ok = newEnforcements[k]; ok {
				if len(strValue) > 0 {
					tflog.Debug(ctx, fmt.Sprintf("add enforcements: %s=%s", k, strValue))
					var re = enterprise.NewRoleEnforcement(roleId, k)
					re.SetValue(strValue)
					enforcements = append(enforcements, re)
				}
			}
		}
	}

	keys = api.NewSet[string]()
	keys.Union(newKeys.ToArray())
	keys.Intersect(oldKeys.ToArray())
	var strValueOrig string
	if len(keys) > 0 { // update
		for _, k := range keys.ToArray() {
			if strValue, ok = newEnforcements[k]; ok {
				if strValueOrig, ok = origEnforcements[k]; ok {
					if strValue != strValueOrig {
						tflog.Debug(ctx, fmt.Sprintf("update enforcements: %s=%s", k, strValue))
						var re = enterprise.NewRoleEnforcement(roleId, k)
						re.SetValue(strValue)
						enforcements = append(enforcements, re)
					}
				}
			}
		}
	}

	if len(enforcements) > 0 {
		var errs = rmr.management.ModifyRoleEnforcements(enforcements)
		for _, er := range errs {
			diags.AddError(
				fmt.Sprintf("Role ID \"%d\" enforcements", roleId),
				fmt.Sprintf("Error: %s", er),
			)
		}
	}
	return
}

func (rmr *roleEnforcementsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan roleEnforcementsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(rmr.applyEnforcements(ctx, plan)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (rmr *roleEnforcementsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan roleEnforcementsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(rmr.applyEnforcements(ctx, plan)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (rmr *roleEnforcementsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleEnforcementsResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	var plan roleEnforcementsResourceModel
	plan.RoleId = state.RoleId
	resp.Diagnostics.Append(rmr.applyEnforcements(ctx, plan)...)
}
