---
page_title: "Quickstart - descope Provider"
description: |-
  Get started with the Descope Terraform Provider. Create your first project and configure authentication methods in minutes.
---

# Quickstart

This guide walks you through setting up the Descope Terraform Provider and creating your first managed project.

## Prerequisites

Before you begin, make sure you have:

1. **Terraform CLI** – [Install Terraform](https://developer.hashicorp.com/terraform/install)
2. **Descope Pro or Enterprise account** – The Terraform provider requires a paid plan
3. **A Management Key** – Create one in [Company Settings](https://app.descope.com/settings/company) in the Descope console. Set the scope to **All Projects** to allow creating new projects.

## Step 1: Configure the Provider

Create a new directory for your Terraform configuration:

```shell
mkdir descope-terraform && cd descope-terraform
```

Create `main.tf` with the provider declaration:

```hcl
terraform {
  required_providers {
    descope = {
      source = "descope/descope"
      version = "~> 1.0"
    }
  }
}

provider "descope" {}
```

Set your management key as an environment variable to avoid hardcoding it:

```shell
export DESCOPE_MANAGEMENT_KEY="K2..."
```

Run `terraform init` to download the provider:

```shell
terraform init
```

## Step 2: Create a Project

Add a project resource to `main.tf`:

```hcl
resource "descope_project" "myapp" {
  name = "my-app"
  tags = ["development"]
}
```

Preview the changes:

```shell
terraform plan
```

Apply them:

```shell
terraform apply
```

Your new Descope project will appear in the [Descope console](https://app.descope.com) within seconds.

## Step 3: Configure Authentication

Add settings resources for the authentication methods your app will use. This example configures Magic Link and Password:

```hcl
resource "descope_magiclink_settings" "myapp" {
  project_id      = descope_project.myapp.id
  expiration_time = "1 hour"
}

resource "descope_password_settings" "myapp" {
  project_id    = descope_project.myapp.id
  min_length    = 10
  lock          = true
  lock_attempts = 5
}
```

Run `terraform apply` to push the updated configuration.

## Step 4: Add Roles and Permissions

Configure RBAC to control what authenticated users can do in your application. Roles and
permissions are standalone resources that reference the project by ID:

```hcl
resource "descope_permission" "read_data" {
  project_id = descope_project.myapp.id
  name       = "read:data"
}

resource "descope_permission" "write_data" {
  project_id = descope_project.myapp.id
  name       = "write:data"
}

resource "descope_role" "viewer" {
  project_id  = descope_project.myapp.id
  name        = "viewer"
  permissions = [descope_permission.read_data.name]
}

resource "descope_role" "editor" {
  project_id  = descope_project.myapp.id
  name        = "editor"
  permissions = [descope_permission.read_data.name, descope_permission.write_data.name]
}
```

## Step 5: Load Authentication Flows

If you've designed custom flows in the Descope console, you can export them and load them via Terraform:

1. In the Descope console, go to **Authentication Flows**
2. Open the flow you want to manage, click the export button, and save the JSON file (e.g., `flows/sign-up-or-in.json`)
3. Reference the file in your configuration:

```hcl
resource "descope_flow" "sign_up_or_in" {
  project_id = descope_project.myapp.id
  flow_id    = "sign-up-or-in"
  data       = file("${path.module}/flows/sign-up-or-in.json")
}
```

## Next Steps

- Explore the resources the provider offers in the [provider reference](../index)
- Browse the [Descope documentation](https://docs.descope.com) for more on authentication concepts
