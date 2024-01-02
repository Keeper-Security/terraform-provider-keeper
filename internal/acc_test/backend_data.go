package acc_test

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
	"github.com/keeper-security/keeper-sdk-golang/vault"
	"strconv"
	"terraform-provider-keeper/internal/model"
)

var _ enterprise.IEnterpriseManagement = new(testingManagement)

//
//	Enterprise ID:    1234    (5299989643264 base enterprise ID)
// 	Enterprise name:  Keeper TF
//
// 	[Keeper TF]      (5299989643266)
//	 +-- [Subnode]   (5299989643274)
//	 +-- Role(s)
//	 |   +--Keeper Administrator (5299989643267)
//   +-- User(s)
//   |   +-- user@company.com    (5299989643284)
//   |   +-- pending_user@company.com    (5299989643285)
//   +-- Team(s)
//   |   +-- Everyone            (MWaZlKLGNa585bX6sCui3g)

func NewTestingManagement() enterprise.IEnterpriseManagement {
	var lm = &testingManagement{
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
		recordTypes:      make(testEntity[vault.IRecordType, string]),
	}
	var rn = enterprise.NewNode(int64(lm.enterpriseId)<<32 + 2)
	rn.SetName(lm.enterpriseInfo.EnterpriseName())
	lm.nodes[rn.NodeId()] = rn
	lm.rootNode = rn

	var sn = enterprise.NewNode(int64(lm.enterpriseId)<<32 + 10)
	sn.SetName("Subnode")
	sn.SetParentId(rn.NodeId())
	sn.SetRestrictVisibility(true)
	lm.nodes[sn.NodeId()] = sn

	var ra = enterprise.NewRole(int64(lm.enterpriseId)<<32 + 3)
	ra.SetName("Keeper Administrator")
	ra.SetNodeId(rn.NodeId())
	ra.SetVisibleBelow(true)
	ra.SetNewUserInherit(false)
	ra.SetKeyType("encrypted_by_data_key")
	lm.roles[ra.RoleId()] = ra

	var re = enterprise.NewRoleEnforcement(ra.RoleId(), "master_password_minimum_length")
	re.SetValue(strconv.Itoa(20))
	lm.roleEnforcements.addLink(re)
	re = enterprise.NewRoleEnforcement(ra.RoleId(), "restrict_email_change")
	re.SetValue("true")
	lm.roleEnforcements.addLink(re)

	var mn = enterprise.NewManagedNode(ra.RoleId(), rn.NodeId())
	mn.SetCascadeNodeManagement(true)
	lm.managedNodes.addLink(mn)

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
	var np = enterprise.NewRolePrivilege(mn.RoleId(), mn.ManagedNodeId())
	pds.ToKeeper(np)
	lm.rolePrivileges.addLink(np)

	var teamUid = "MWaZlKLGNa585bX6sCui3g"
	var t = enterprise.NewTeam(teamUid)
	t.SetName("Everyone")
	t.SetNodeId(rn.NodeId())
	lm.teams[teamUid] = t

	var eUid = int64(lm.enterpriseId)<<32 + 20
	var us = enterprise.NewUser(eUid, "user@company.com", "active")
	us.SetNodeId(rn.NodeId())
	us.SetFullName("Keeper Admin")
	lm.users[us.EnterpriseUserId()] = us

	var teamUser = enterprise.NewTeamUser(teamUid, eUid)
	lm.teamUsers.addLink(teamUser)

	var roleUser = enterprise.NewRoleUser(ra.RoleId(), eUid)
	lm.roleUsers.addLink(roleUser)

	eUid = int64(lm.enterpriseId)<<32 + 21
	us = enterprise.NewUser(eUid, "pending_user@company.com", "inactive")
	us.SetNodeId(rn.NodeId())
	us.SetFullName("Invited User")
	lm.users[us.EnterpriseUserId()] = us

	var renf = enterprise.NewRoleEnforcement(ra.RoleId(), "restrict_record_types")
	renf.SetValue(`{"std":[11,14],"ent":[]}`)
	lm.roleEnforcements.addLink(renf)

	var err error
	var irt vault.IRecordType
	var rt = &storageRecordType{
		id:      14,
		scope:   vault.RecordTypeScope_Standard,
		content: `{"$id":"address","categories":["address"],"description":"Address template","fields":[{"$ref":"address"},{"$ref":"fileRef"}]}`,
	}
	if irt, err = vault.ParseRecordType(rt); err == nil {
		lm.recordTypes["address"] = irt
	}
	rt = &storageRecordType{
		id:      1,
		scope:   vault.RecordTypeScope_Standard,
		content: `{"$id":"login","categories":["login"],"description":"Login template","fields":[{"$ref":"passkey"},{"$ref":"login"},{"$ref":"password"},{"$ref":"url"},{"$ref":"fileRef"},{"$ref":"oneTimeCode"}]}`,
	}
	if irt, err = vault.ParseRecordType(rt); err == nil {
		lm.recordTypes["login"] = irt
	}
	rt = &storageRecordType{
		id:      11,
		scope:   vault.RecordTypeScope_Standard,
		content: `{"$id":"bankAccount","description":"Bank account template","fields":[{"$ref":"bankAccount","required":true},{"$ref":"name"},{"$ref":"login"},{"$ref":"password"},{"$ref":"url"},{"$ref":"cardRef"},{"$ref":"fileRef"},{"$ref":"oneTimeCode"}]}`,
	}
	if irt, err = vault.ParseRecordType(rt); err == nil {
		lm.recordTypes["bankAccount"] = irt
	}

	rt = &storageRecordType{
		id:      18,
		scope:   vault.RecordTypeScope_Standard,
		content: `{"$id":"bankCard","categories":["payment"],"description":"Bank card template","fields":[{"$ref":"paymentCard"},{"$ref":"text","label":"cardholderName"},{"$ref":"pinCode"},{"$ref":"addressRef"},{"$ref":"fileRef"}]}`,
	}
	if irt, err = vault.ParseRecordType(rt); err == nil {
		lm.recordTypes["bankCard"] = irt
	}
	rt = &storageRecordType{
		id:      20,
		scope:   vault.RecordTypeScope_Standard,
		content: `{"$id":"contact","categories":["address"],"description":"Contact template","fields":[{"$ref":"name","required":true},{"$ref":"text","label":"company"},{"$ref":"email"},{"$ref":"phone"},{"$ref":"addressRef"},{"$ref":"fileRef"}]}`,
	}
	if irt, err = vault.ParseRecordType(rt); err == nil {
		lm.recordTypes["contact"] = irt
	}
	return lm
}

type storageRecordType struct {
	id      int64
	scope   vault.RecordTypeScope
	content string
}

func (srt *storageRecordType) Id() int64 {
	return srt.id
}
func (srt *storageRecordType) Scope() int32 {
	return int32(srt.scope)
}
func (srt *storageRecordType) Content() string {
	return srt.content
}
func (srt *storageRecordType) Uid() int64 {
	return srt.Id()
}
