terraform {
  required_providers {
    kepr = {
      source = "keepersecurity.com/Keeper-Security/kepr"
    }
  }
}

provider "kepr" {
   // config_type = "commander"
}







/*
resource "keeper-enterprise_team" "Team1" {
  name=
  team_uid =
  restrict_view
}
resource "keeper-enterprise_user" "User1" {
  email
  full_name = "lll"
  active
  lock
}

resource "keeper-enterprise_users_team" "Users" {
  team = "abc"
  users = data.dfdfdfdf
}

data "keeper-enterprise_users_teams" "Users" { List [TeamName <-> User]
  teams = [a,c,d]
  users = [1,2,3,4]
  List[]
}

resource "keeper-enterprise_users_teams" "Users" {
  exists = false
  users = [1,2,3,4]
  teams = [e]   script
}

resource "keeper-enterprise_user_team" "Administrators" {
  users = [1,2,3,10]
  teams = [e,f,g]   script
}

resource "keeper-enterprise_users_roles" "Administrators" {
  users = [1,2,3,10]
  roles = [e,f,g]   script
}

resource "keeper-enterprise_user_team" "Administrator1" {
  exists = true
  users = [2,6]
  teams = [g,h]   script
}


resource "keeper-enterprise_user" "user" {
  email = var.email
  teams = data.keeper-enterprise_teams.teams.team_uids
}

data "keeper-enterprise_user" "user" {
  email = var.email
}

data "keeper-enterprise_teams" "teams" {  // -> IDs
  names = [a,b,c,d]
}

*/
#data "kepr_nodes" "azure" {
#  subnodes = {
#    node_id = 820338837337
#    cascade = true
#    include_parent = true
#  }
#}

#data "kepr_node" "root" {
#  is_root = true
#  //name = "nnd"
#}

#data "kepr_users" "active" {
#  emails = ["skol@skolupaev.info", "sergey+a2@callpod.com"]
#}

#resource "kepr_team" "team33" {
#    name = "Team33"
#    node_id = 820338837337
#    restrict_share = true
#}

data "kepr_team" "vault" {
  name = "Vault"
  include_users = true
}

#resource "kepr_team_users" "vault_users" {
#  team_uid = resource.vault.team_uid
#}
#data "kepr_login_enforcements" "user_login" {
#  minimum = 12
#  upper = 2
#}
#
#resource "kepr_role"  "aaa" {
#  login_enforments = {
#    minimum = 12
#  }
#}

 output "example" {
   value = data.kepr_team.vault
 }