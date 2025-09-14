![terraform-provider-keeper Header](https://github.com/user-attachments/assets/fd58ea10-ae8f-4642-902c-65006ace6425)

# Keeper Terraform Provider (Administration)

[![Terraform Registry](https://img.shields.io/badge/registry-keeper-blue)](https://registry.terraform.io/providers/Keeper-Security/keeper/latest/docs)

Manage Keeper administration as code with Terraform: users, teams (groups), roles, and role enforcements/policies.

---

## Features

- **Users**: Look up and manage user entries.
- **Teams / Groups**: Work with teams and memberships.
- **Roles & Enforcements**: Define and apply role enforcements/policies.

> Full docs are in the Terraform Registry for the `keeper` provider.

---

## Requirements

- Terraform (latest stable recommended)
- Keeper admin access (sufficient privileges to manage users/teams/roles)
- A Keeper provider configuration file (JSON) accessible via `config_path` (see below). The Registry’s provider page shows `config_path` usage. 

---

## Installation

Declare the provider in your Terraform configuration:

```hcl
terraform {
  required_providers {
    keeper = {
      source  = "Keeper-Security/keeper"
      # pin as needed, e.g. version = "~> 1.2"
    }
  }
}
