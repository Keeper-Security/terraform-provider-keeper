package acc_test

import (
	"crypto/ecdh"
	"crypto/rsa"
	"github.com/keeper-security/keeper-sdk-golang/api"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
	"github.com/keeper-security/keeper-sdk-golang/storage"
	"github.com/keeper-security/keeper-sdk-golang/vault"
	"sync/atomic"
)

type testEntity[T storage.IUid[K], K storage.Key] map[K]T

func (te testEntity[T, K]) GetAllEntities(cb func(T) bool) {
	for _, v := range te {
		if !cb(v) {
			return
		}
	}
}
func (te testEntity[T, K]) GetEntity(uid K) (t T) {
	t = te[uid]
	return
}
func (te testEntity[T, K]) deleteEntity(uid K) {
	delete(te, uid)
}

type testLink[T storage.IUidLink[KS, KO], KS storage.Key, KO storage.Key] map[enterprise.LinkKey[KS, KO]]T

func (tl testLink[T, KS, KO]) GetLink(subjectUid KS, objectUid KO) (t T) {
	var link = enterprise.LinkKey[KS, KO]{
		V1: subjectUid,
		V2: objectUid,
	}
	t = tl[link]
	return
}
func (tl testLink[T, KS, KO]) GetLinksBySubject(subjectUid KS, cb func(T) bool) {
	for k, v := range tl {
		if k.V1 == subjectUid {
			if !cb(v) {
				return
			}
		}
	}
}
func (tl testLink[T, KS, KO]) GetLinksByObject(objectUid KO, cb func(T) bool) {
	for k, v := range tl {
		if k.V2 == objectUid {
			if !cb(v) {
				return
			}
		}
	}
}
func (tl testLink[T, KS, KO]) GetAllLinks(cb func(T) bool) {
	for _, v := range tl {
		if !cb(v) {
			return
		}
	}
}
func (tl testLink[T, KS, KO]) addLink(t T) {
	var link = enterprise.LinkKey[KS, KO]{
		V1: t.SubjectUid(),
		V2: t.ObjectUid(),
	}
	tl[link] = t
}
func (tl testLink[T, KS, KO]) deleteLink(link enterprise.LinkKey[KS, KO]) {
	delete(tl, link)
}
func (tl testLink[T, KS, KO]) deleteSubject(subjectUid KS) {
	for k := range tl {
		if k.V1 == subjectUid {
			delete(tl, k)
		}
	}
}
func (tl testLink[T, KS, KO]) deleteObject(objectUid KO) {
	for k := range tl {
		if k.V2 == objectUid {
			delete(tl, k)
		}
	}
}

func newTestEnterpriseInfo() *testEnterpriseInfo {
	var err error
	var rsaKey *rsa.PrivateKey
	if rsaKey, _, err = api.GenerateRsaKey(); err != nil {
		panic(err)
	}

	var ecKey *ecdh.PrivateKey
	if ecKey, _, err = api.GenerateEcKey(); err != nil {
		panic(err)
	}
	return &testEnterpriseInfo{
		treeKey:       api.GenerateAesKey(),
		rsaPrivateKey: rsaKey,
		ecPrivateKey:  ecKey,
	}
}

var _ enterprise.IEnterpriseInfo = new(testEnterpriseInfo)

type testEnterpriseInfo struct {
	treeKey       []byte
	rsaPrivateKey *rsa.PrivateKey
	ecPrivateKey  *ecdh.PrivateKey
}

func (ei *testEnterpriseInfo) EnterpriseName() string {
	return "Keeper TF"
}
func (ei *testEnterpriseInfo) IsDistributor() bool {
	return false
}
func (ei *testEnterpriseInfo) TreeKey() []byte {
	return ei.treeKey
}
func (ei *testEnterpriseInfo) RsaPrivateKey() *rsa.PrivateKey {
	return ei.rsaPrivateKey
}
func (ei *testEnterpriseInfo) EcPrivateKey() *ecdh.PrivateKey {
	return ei.ecPrivateKey
}

type testingManagement struct {
	enterpriseId     int32
	lastId           int32
	enterpriseInfo   *testEnterpriseInfo
	nodes            testEntity[enterprise.INode, int64]
	rootNode         enterprise.INode
	roles            testEntity[enterprise.IRole, int64]
	users            testEntity[enterprise.IUser, int64]
	teams            testEntity[enterprise.ITeam, string]
	teamUsers        testLink[enterprise.ITeamUser, string, int64]
	queuedTeams      testEntity[enterprise.IQueuedTeam, string]
	queuedTeamUsers  testLink[enterprise.IQueuedTeamUser, string, int64]
	roleUsers        testLink[enterprise.IRoleUser, int64, int64]
	roleTeams        testLink[enterprise.IRoleTeam, int64, string]
	managedNodes     testLink[enterprise.IManagedNode, int64, int64]
	rolePrivileges   testLink[enterprise.IRolePrivilege, int64, int64]
	roleEnforcements testLink[enterprise.IRoleEnforcement, int64, string]
	licenses         testEntity[enterprise.ILicense, int64]
	userAliases      testLink[enterprise.IUserAlias, int64, string]
	ssoServices      testEntity[enterprise.ISsoService, int64]
	bridges          testEntity[enterprise.IBridge, int64]
	scims            testEntity[enterprise.IScim, int64]
	managedCompanies testEntity[enterprise.IManagedCompany, int64]
	recordTypes      testEntity[vault.IRecordType, string]
}

func (lm *testingManagement) EnterpriseData() enterprise.IEnterpriseData {
	return lm
}
func (lm *testingManagement) EnterpriseInfo() enterprise.IEnterpriseInfo {
	return lm.enterpriseInfo
}
func (lm *testingManagement) Nodes() enterprise.IEnterpriseEntity[enterprise.INode, int64] {
	return lm.nodes
}
func (lm *testingManagement) RootNode() enterprise.INode {
	return lm.rootNode
}
func (lm *testingManagement) Roles() enterprise.IEnterpriseEntity[enterprise.IRole, int64] {
	return lm.roles
}
func (lm *testingManagement) Users() enterprise.IEnterpriseEntity[enterprise.IUser, int64] {
	return lm.users
}
func (lm *testingManagement) Teams() enterprise.IEnterpriseEntity[enterprise.ITeam, string] {
	return lm.teams
}
func (lm *testingManagement) TeamUsers() enterprise.IEnterpriseLink[enterprise.ITeamUser, string, int64] {
	return lm.teamUsers
}
func (lm *testingManagement) QueuedTeams() enterprise.IEnterpriseEntity[enterprise.IQueuedTeam, string] {
	return lm.queuedTeams
}
func (lm *testingManagement) QueuedTeamUsers() enterprise.IEnterpriseLink[enterprise.IQueuedTeamUser, string, int64] {
	return lm.queuedTeamUsers
}
func (lm *testingManagement) RoleUsers() enterprise.IEnterpriseLink[enterprise.IRoleUser, int64, int64] {
	return lm.roleUsers
}
func (lm *testingManagement) RoleTeams() enterprise.IEnterpriseLink[enterprise.IRoleTeam, int64, string] {
	return lm.roleTeams
}
func (lm *testingManagement) ManagedNodes() enterprise.IEnterpriseLink[enterprise.IManagedNode, int64, int64] {
	return lm.managedNodes
}
func (lm *testingManagement) RolePrivileges() enterprise.IEnterpriseLink[enterprise.IRolePrivilege, int64, int64] {
	return lm.rolePrivileges
}
func (lm *testingManagement) RoleEnforcements() enterprise.IEnterpriseLink[enterprise.IRoleEnforcement, int64, string] {
	return lm.roleEnforcements
}
func (lm *testingManagement) Licenses() enterprise.IEnterpriseEntity[enterprise.ILicense, int64] {
	return lm.licenses
}
func (lm *testingManagement) UserAliases() enterprise.IEnterpriseLink[enterprise.IUserAlias, int64, string] {
	return lm.userAliases
}
func (lm *testingManagement) SsoServices() enterprise.IEnterpriseEntity[enterprise.ISsoService, int64] {
	return lm.ssoServices
}
func (lm *testingManagement) Bridges() enterprise.IEnterpriseEntity[enterprise.IBridge, int64] {
	return lm.bridges
}
func (lm *testingManagement) Scims() enterprise.IEnterpriseEntity[enterprise.IScim, int64] {
	return lm.scims
}
func (lm *testingManagement) ManagedCompanies() enterprise.IEnterpriseEntity[enterprise.IManagedCompany, int64] {
	return lm.managedCompanies
}
func (lm *testingManagement) RecordTypes() enterprise.IEnterpriseEntity[vault.IRecordType, string] {
	return lm.recordTypes
}

func (lm *testingManagement) GetEnterpriseId() (id int64, err error) {
	id = int64(lm.enterpriseId)<<32 + int64(atomic.AddInt32(&lm.lastId, 1))
	return
}
func (lm *testingManagement) ModifyNodes(nodesToAdd []enterprise.INode, nodesToUpdate []enterprise.INode, nodesToDelete []int64) (errs []error) {
	var ok bool
	for _, node := range nodesToUpdate {
		var nodeId = node.NodeId()
		if (nodeId&0xff == 2) && node.ParentId() > 0 {
			errs = append(errs, api.NewKeeperApiError("bad_inputs_parent_id", "can't move root"))
		} else if node.NodeId() == node.ParentId() {
			errs = append(errs, api.NewKeeperApiError("bad_inputs_parent_id", "parent node is self"))
		} else if _, ok = lm.nodes[nodeId]; ok {
			if _, ok = lm.nodes[node.ParentId()]; ok {
				lm.nodes[node.NodeId()] = node
			} else {
				errs = append(errs, api.NewKeeperApiError("doesnt_exist", "parent node"))
			}
		} else {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "node id"))
		}
	}

	for _, node := range nodesToAdd {
		var nodeId = node.NodeId()
		if nodeId>>32 != int64(lm.enterpriseId) {
			errs = append(errs, api.NewKeeperApiError("bad_inputs_node_id", "out of range"))
		} else if _, ok = lm.nodes[nodeId]; !ok {
			if _, ok = lm.nodes[node.ParentId()]; ok {
				lm.nodes[node.NodeId()] = node
			} else {
				errs = append(errs, api.NewKeeperApiError("doesnt_exist", "parent node"))
			}
		} else {
			errs = append(errs, api.NewKeeperApiError("bad_inputs_node_id", "already exists"))
		}
	}

	var parents = api.NewSet[int64]()
	lm.nodes.GetAllEntities(func(r enterprise.INode) bool {
		if r.ParentId() > 0 {
			parents.Add(r.ParentId())
		}
		return true
	})
	lm.roles.GetAllEntities(func(r enterprise.IRole) bool {
		parents.Add(r.NodeId())
		return true
	})
	lm.users.GetAllEntities(func(u enterprise.IUser) bool {
		parents.Add(u.NodeId())
		return true
	})
	for _, nodeId := range nodesToDelete {
		if parents.Has(nodeId) {
			errs = append(errs, api.NewKeeperApiError("node_not_empty", "delete node"))
		} else if _, ok = lm.nodes[nodeId]; ok {
			lm.nodes.deleteEntity(nodeId)
			parents.Delete(nodeId)
		} else {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "existing node"))
		}
	}
	return
}
func (lm *testingManagement) ModifyRoles(rolesToAdd []enterprise.IRole, rolesToUpdate []enterprise.IRole, rolesToDelete []int64) (errs []error) {
	var ok bool
	for _, role := range rolesToUpdate {
		if _, ok = lm.nodes[role.NodeId()]; !ok {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "parent node"))
		} else if _, ok = lm.roles[role.RoleId()]; !ok {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "existing role"))
		} else {
			if role.NewUserInherit() {
				var isAdmin = false
				lm.managedNodes.GetLinksBySubject(role.RoleId(), func(mn enterprise.IManagedNode) bool {
					isAdmin = true
					return false
				})
				if isAdmin {
					errs = append(errs, api.NewKeeperApiError("cant_add_inherit", "admin role cannot be inherited"))
					continue
				}
			}
			lm.roles[role.RoleId()] = role
		}
	}

	for _, role := range rolesToAdd {
		var roleId = role.RoleId()
		if roleId>>32 != int64(lm.enterpriseId) {
			errs = append(errs, api.NewKeeperApiError("bad_inputs_role_id", "out of range"))
		} else if _, ok = lm.roles[roleId]; !ok {
			if _, ok = lm.nodes[role.NodeId()]; ok {
				lm.roles[roleId] = role
			} else {
				errs = append(errs, api.NewKeeperApiError("doesnt_exist", "parent node"))
			}
		} else {
			errs = append(errs, api.NewKeeperApiError("exists", "already exists"))
		}
	}

	for _, roleId := range rolesToDelete {
		if _, ok = lm.roles[roleId]; ok {
			lm.roles.deleteEntity(roleId)
			lm.roleUsers.deleteSubject(roleId)
			lm.roleEnforcements.deleteSubject(roleId)
			lm.managedNodes.deleteSubject(roleId)
			lm.rolePrivileges.deleteSubject(roleId)
		} else {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "existing role"))
		}
	}

	return
}
func (lm *testingManagement) ModifyTeams(teamsToAdd []enterprise.ITeam, teamsToUpdate []enterprise.ITeam, teamsToDelete []string) (errs []error) {
	var ok bool
	for _, team := range teamsToUpdate {
		var nodeId = team.NodeId()
		if _, ok = lm.nodes[nodeId]; !ok {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "parent node"))
		} else if _, ok = lm.teams[team.TeamUid()]; !ok {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "existing team"))
		} else {
			lm.teams[team.TeamUid()] = team
		}
	}

	for _, team := range teamsToAdd {
		var teamUid = team.TeamUid()
		if _, ok = lm.teams[teamUid]; !ok {
			if _, ok = lm.nodes[team.NodeId()]; ok {
				lm.teams[teamUid] = team
			} else {
				errs = append(errs, api.NewKeeperApiError("doesnt_exist", "parent node"))
			}
		} else {
			errs = append(errs, api.NewKeeperApiError("exists", "already exists"))
		}
	}

	for _, teamUid := range teamsToDelete {
		if _, ok = lm.teams[teamUid]; ok {
			lm.teams.deleteEntity(teamUid)
			lm.teamUsers.deleteSubject(teamUid)
			lm.queuedTeamUsers.deleteSubject(teamUid)
			lm.roleTeams.deleteObject(teamUid)
		} else {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "existing team"))
		}
	}

	return
}

func (lm *testingManagement) ModifyTeamUsers(teamUsersToAdd []enterprise.ITeamUser, teamUsersToRemove []enterprise.ITeamUser) (errs []error) {
	for _, teamUser := range teamUsersToAdd {
		var teamUid = teamUser.TeamUid()
		var eUserId = teamUser.EnterpriseUserId()
		var l = lm.teamUsers.GetLink(teamUid, eUserId)
		if l == nil {
			if team := lm.teams.GetEntity(teamUid); team != nil {
				if user := lm.users.GetEntity(eUserId); user != nil {
					if user.Status() == "active" {
						lm.teamUsers.addLink(teamUser)
					} else {
						errs = append(errs, api.NewKeeperApiError("cant_be_pending", "user cannot be pending"))
					}
				} else {
					errs = append(errs, api.NewKeeperApiError("doesnt_exist", "user does not exist"))
				}
			} else {
				errs = append(errs, api.NewKeeperApiError("doesnt_exist", "team does not exist"))
			}
		} else {
			errs = append(errs, api.NewKeeperApiError("exists", "already exists"))
		}
	}

	for _, teamUser := range teamUsersToRemove {
		var teamUid = teamUser.TeamUid()
		var eUserId = teamUser.EnterpriseUserId()
		var l = lm.teamUsers.GetLink(teamUid, eUserId)
		if l != nil {
			var link = enterprise.LinkKey[string, int64]{
				V1: teamUid,
				V2: eUserId,
			}
			lm.teamUsers.deleteLink(link)
		} else {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "team membership does not exist"))
		}
	}
	return
}
func (lm *testingManagement) ModifyRoleUsers(roleUsersToAdd []enterprise.IRoleUser, roleUsersToRemove []enterprise.IRoleUser) (errs []error) {
	for _, roleUser := range roleUsersToAdd {
		var roleId = roleUser.RoleId()
		var eUserId = roleUser.EnterpriseUserId()
		var l = lm.roleUsers.GetLink(roleId, eUserId)
		if l == nil {
			if role := lm.roles.GetEntity(roleId); role != nil {
				var isAdmin = false
				lm.managedNodes.GetLinksBySubject(roleId, func(x enterprise.IManagedNode) bool {
					isAdmin = true
					return false
				})
				if user := lm.users.GetEntity(eUserId); user != nil {
					if isAdmin && user.Status() == "active" {
						lm.roleUsers.addLink(roleUser)
					} else {
						errs = append(errs, api.NewKeeperApiError("cant_be_pending", "admin role: user cannot be pending"))
					}
				} else {
					errs = append(errs, api.NewKeeperApiError("doesnt_exist", "user does not exist"))
				}
			} else {
				errs = append(errs, api.NewKeeperApiError("doesnt_exist", "role does not exist"))
			}
		} else {
			errs = append(errs, api.NewKeeperApiError("exists", "already exists"))
		}
	}

	for _, roleUser := range roleUsersToRemove {
		var roleId = roleUser.RoleId()
		var eUserId = roleUser.EnterpriseUserId()
		var l = lm.roleUsers.GetLink(roleId, eUserId)
		if l != nil {
			var link = enterprise.LinkKey[int64, int64]{
				V1: roleId,
				V2: eUserId,
			}
			lm.roleUsers.deleteLink(link)
		} else {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "role membership does not exist"))
		}
	}
	return
}

func (lm *testingManagement) ModifyRoleTeams(roleTeamsToAdd []enterprise.IRoleTeam, roleTeamsToRemove []enterprise.IRoleTeam) (errs []error) {
	for _, roleTeam := range roleTeamsToAdd {
		var roleId = roleTeam.RoleId()
		var teamUid = roleTeam.TeamUid()
		var l = lm.roleTeams.GetLink(roleId, teamUid)
		if l == nil {
			if role := lm.roles.GetEntity(roleId); role != nil {
				var isAdmin = false
				lm.managedNodes.GetLinksBySubject(roleId, func(x enterprise.IManagedNode) bool {
					isAdmin = true
					return false
				})
				if !isAdmin {
					if team := lm.teams.GetEntity(teamUid); team != nil {
						lm.roleTeams.addLink(roleTeam)
					} else {
						errs = append(errs, api.NewKeeperApiError("doesnt_exist", "team does not exist"))
					}

				} else {
					errs = append(errs, api.NewKeeperApiError("cant_add_admin_role", "admin role: team cannot be added"))
				}
			} else {
				errs = append(errs, api.NewKeeperApiError("doesnt_exist", "role does not exist"))
			}
		} else {
			errs = append(errs, api.NewKeeperApiError("exists", "already exists"))
		}
	}

	for _, roleTeam := range roleTeamsToRemove {
		var roleId = roleTeam.RoleId()
		var teamUid = roleTeam.TeamUid()
		var l = lm.roleTeams.GetLink(roleId, teamUid)
		if l != nil {
			var link = enterprise.LinkKey[int64, string]{
				V1: roleId,
				V2: teamUid,
			}
			lm.roleTeams.deleteLink(link)
		} else {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "role-team does not exist"))
		}
	}
	return
}

func (lm *testingManagement) ModifyManagedNodes(managedNodesToAdd []enterprise.IManagedNode, managedNodesToUpdate []enterprise.IManagedNode, managedNodesToRemove []enterprise.IManagedNode) (errs []error) {
	var ok bool
	var link enterprise.LinkKey[int64, int64]
	for _, managedNode := range managedNodesToUpdate {
		link.V1 = managedNode.RoleId()
		link.V2 = managedNode.ManagedNodeId()
		if _, ok = lm.managedNodes[link]; !ok {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "managed node"))
		} else {
			lm.managedNodes[link] = managedNode
		}
	}

	for _, managedNode := range managedNodesToAdd {
		link.V1 = managedNode.RoleId()
		link.V2 = managedNode.ManagedNodeId()
		if _, ok = lm.managedNodes[link]; !ok {
			if _, ok = lm.roles[managedNode.RoleId()]; ok {
				if _, ok = lm.nodes[managedNode.ManagedNodeId()]; ok {
					lm.managedNodes[link] = managedNode
				} else {
					errs = append(errs, api.NewKeeperApiError("doesnt_exist", "node"))
				}
			} else {
				errs = append(errs, api.NewKeeperApiError("doesnt_exist", "role"))
			}
		} else {
			errs = append(errs, api.NewKeeperApiError("exists", "already exists"))
		}
	}

	for _, managedNode := range managedNodesToRemove {
		link.V1 = managedNode.RoleId()
		link.V2 = managedNode.ManagedNodeId()
		if _, ok = lm.managedNodes[link]; ok {
			lm.managedNodes.deleteLink(link)
			lm.rolePrivileges.deleteLink(link)
		} else {
			errs = append(errs, api.NewKeeperApiError("doesnt_exist", "existing managed node"))
		}
	}

	return
}

func (lm *testingManagement) ModifyRolePrivileges(privileges []enterprise.IRolePrivilege) (errs []error) {
	for _, e := range privileges {
		lm.rolePrivileges.addLink(e)
	}
	return
}

func (lm *testingManagement) ModifyRoleEnforcements(enforcementsToSet []enterprise.IRoleEnforcement) (err []error) {
	for _, e := range enforcementsToSet {
		if len(e.Value()) > 0 {
			lm.roleEnforcements.addLink(e)
		} else {
			lm.roleEnforcements.deleteLink(enterprise.LinkKey[int64, string]{
				V1: e.RoleId(),
				V2: e.EnforcementType(),
			})
		}
	}
	return
}

func (lm *testingManagement) Commit() (errs []error) {
	return
}
