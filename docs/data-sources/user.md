---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "keeper_user Data Source - terraform-provider-keeper"
subcategory: ""
description: |-
  
---

# keeper_user (Data Source)



## Example Usage

```terraform
data "keeper_user" "example_by_username" {
  username      = "user@company.com"
  include_roles = true
  include_teams = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `enterprise_user_id` (Number) Enterprise User ID
- `include_roles` (Boolean) Include user roles
- `include_teams` (Boolean) Include user teams
- `username` (String) User email

### Read-Only

- `account_share_expiration` (Number) Account Share deadline: Timestamp in millis
- `full_name` (String) User Full Name
- `job_title` (String) User Job Title
- `node_id` (Number) User Node ID
- `roles` (Attributes List) (see [below for nested schema](#nestedatt--roles))
- `status` (String) User Status: active | invited | locked | blocked | disabled
- `teams` (Attributes List) (see [below for nested schema](#nestedatt--teams))
- `tfa_enabled` (Boolean) TFA Enabled flag

<a id="nestedatt--roles"></a>
### Nested Schema for `roles`

Read-Only:

- `is_admin` (Boolean) Is Administrative Role
- `name` (String) Role Name
- `node_id` (Number) Role Node ID
- `role_id` (Number) Role ID


<a id="nestedatt--teams"></a>
### Nested Schema for `teams`

Read-Only:

- `name` (String) Team Name
- `team_uid` (String) Team UID
