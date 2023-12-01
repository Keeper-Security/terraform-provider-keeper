package test

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"terraform-provider-kepr/internal/model"
)

var _ enterprise.IEnterpriseManagement = new(loopbackManagement)

func NewloopbackManagement() (lm *loopbackManagement) {
	lm = &loopbackManagement{
		enterpriseId:     1234,
		lastId:           1000,
		enterpriseInfo:   newTestEnterpriseInfo(),
		nodes:            make(testEntity[enterprise.INode, int64]),
		roles:            make(testEntity[enterprise.IRole, int64]),
		users:            make(testEntity[enterprise.IUser, int64]),
		teams:            make(testEntity[enterprise.ITeam, string]),
		teamUsers:        make(testLink[enterprise.ITeamUser, string, int64]),
		queuedTeams:      make(testEntity[enterprise.IQueuedTeam, string]),
		queuedTeamUsers:  make(testLink[enterprise.IQueuedTeamUser, string, int64]),
		roleUsers:        make(testLink[enterprise.IRoleUser, int64, int64]),
		roleTeams:        make(testLink[enterprise.IRoleTeam, int64, string]),
		managedNodes:     make(testLink[enterprise.IManagedNode, int64, int64]),
		rolePrivileges:   make(testLink[enterprise.IRolePrivilege, int64, int64]),
		roleEnforcements: make(testLink[enterprise.IRoleEnforcement, int64, string]),
		licenses:         make(testEntity[enterprise.ILicense, int64]),
		userAliases:      make(testLink[enterprise.IUserAlias, int64, string]),
		ssoServices:      make(testEntity[enterprise.ISsoService, int64]),
		bridges:          make(testEntity[enterprise.IBridge, int64]),
		scims:            make(testEntity[enterprise.IScim, int64]),
		managedCompanies: make(testEntity[enterprise.IManagedCompany, int64]),
	}
	var rn = enterprise.NewNode(int64(lm.enterpriseId)<<32 + 2)
	rn.SetName(lm.enterpriseInfo.EnterpriseName())
	lm.nodes[rn.NodeId()] = rn
	lm.rootNode = rn

	var ra = enterprise.NewRole(int64(lm.enterpriseId)<<32 + 3)
	ra.SetName("Keeper Administrator")
	ra.SetNodeId(rn.NodeId())
	ra.SetVisibleBelow(true)
	ra.SetNewUserInherit(false)
	ra.SetKeyType("encrypted_by_data_key")
	lm.roles[ra.RoleId()] = ra

	var mn = enterprise.NewManagedNode(ra.RoleId(), rn.NodeId())
	mn.SetCascadeNodeManagement(true)
	lm.managedNodes.addLink(mn)

	var np = enterprise.NewRolePrivilege(mn.RoleId(), mn.ManagedNodeId())
	var pds = &model.PrivilegeDataSourceModel{
		ManageNodes:          types.BoolValue(true),
		ManageUsers:          types.BoolValue(true),
		ManageRoles:          types.BoolValue(true),
		ManageTeams:          types.BoolValue(true),
		ManageReports:        types.BoolValue(true),
		ManageSso:            types.BoolValue(true),
		DeviceApproval:       types.BoolValue(true),
		ManageRecordTypes:    types.BoolValue(true),
		ShareAdmin:           types.BoolValue(true),
		RunComplianceReports: types.BoolValue(true),
		ManageCompanies:      types.BoolValue(true),
		TransferAccount:      types.BoolValue(false),
	}
	for _, p := range pds.ToKeeper() {
		np.SetPrivilege(p)
	}
	lm.rolePrivileges.addLink(np)
	var err error
	var eUid int64
	if eUid, err = lm.GetEnterpriseId(); err != nil {
		panic(err)
	}
	var us = enterprise.NewUser(eUid, "user@company.com")
	us.SetNodeId(rn.NodeId())
	us.SetFullName("Keeper Admin")
	lm.users[us.EnterpriseUserId()] = us
	return
}
