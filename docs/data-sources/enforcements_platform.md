---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "keeper_enforcements_platform Data Source - terraform-provider-keeper"
subcategory: ""
description: |-
  
---

# keeper_enforcements_platform (Data Source)



## Example Usage

```terraform
data "keeper_enforcements_platform" "example" {
  restrict_commander_access    = true
  restrict_extensions_access   = true
  restrict_chat_desktop_access = true
  restrict_chat_mobile_access  = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `restrict_chat_desktop_access` (Boolean) Restrict access to Keeper Chat for desktop
- `restrict_chat_mobile_access` (Boolean) Restrict access to Keeper Chat for mobile
- `restrict_commander_access` (Boolean) Restrict access to Keeper Commander
- `restrict_desktop_access` (Boolean) Restrict access to Keeper Vault for desktop
- `restrict_desktop_mac_access` (Boolean) Restrict access to Keeper Desktop for MacOS
- `restrict_desktop_win_access` (Boolean) Restrict access to Keeper Desktop for Windows
- `restrict_extensions_access` (Boolean) Restrict access to Keeper browser extensions
- `restrict_mobile_access` (Boolean) Restrict access to Keeper Vault for mobile
- `restrict_mobile_android_access` (Boolean) Restrict access to Keeper Vault for Android
- `restrict_mobile_ios_access` (Boolean) Restrict access to Keeper Vault for iOS
- `restrict_mobile_windows_phone_access` (Boolean) Restrict access to Keeper Vault for Windows mobile
- `restrict_web_vault_access` (Boolean) Restrict access to Keeper Vault for web