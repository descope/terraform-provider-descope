---
page_title: "Microsoft Sentinel Audit Stream - descope Provider"
description: |-
  Create the Azure resources that the Microsoft Sentinel audit stream connector writes to.
---

# Microsoft Sentinel Audit Stream

The Microsoft Sentinel connector streams Descope audit events to a Log Analytics workspace through the Azure Monitor
Logs Ingestion API. This guide provides a Terraform module that creates the Azure side of that setup:

- The `DescopeAudit_CL` custom table in the workspace
- A `Direct` data collection rule that routes the `Custom-DescopeAudit_CL` stream into the table
- An Entra app registration, service principal and client secret that Descope authenticates as
- A Monitoring Metrics Publisher role assignment for the service principal on the data collection rule

The module uses the `azurerm`, `azuread` and `azapi` providers. The data collection rule is created with `azapi`
because `azurerm_monitor_data_collection_rule` does not expose the logs ingestion endpoint of a `Direct` rule.

## Module

Save the following files in a `descope-sentinel` directory.

`variables.tf`:

```terraform
variable "resource_group_name" {
  type        = string
  description = "Resource group of the Log Analytics workspace that Microsoft Sentinel is enabled on."
}

variable "log_analytics_workspace_id" {
  type        = string
  description = "Resource ID of the Log Analytics workspace that Microsoft Sentinel is enabled on."
}

variable "location" {
  type        = string
  description = "Azure region of the Log Analytics workspace."
}

variable "application_display_name" {
  type        = string
  description = "Display name of the Entra app registration that Descope authenticates as."
  default     = "Descope Audit Stream"
}
```

`main.tf`:

```terraform
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.55.0"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = ">= 3.0"
    }
    azapi = {
      source  = "azure/azapi"
      version = ">= 2.0"
    }
  }
}

locals {
  table_name  = "DescopeAudit_CL"
  stream_name = "Custom-DescopeAudit_CL"
  columns = [
    { name = "TimeGenerated", type = "dateTime" },
    { name = "AuditId", type = "string" },
    { name = "Action", type = "string" },
    { name = "AuditType", type = "string" },
    { name = "ProjectId", type = "string" },
    { name = "UserId", type = "string" },
    { name = "ActorId", type = "string" },
    { name = "LoginIds", type = "dynamic" },
    { name = "Tenants", type = "dynamic" },
    { name = "Method", type = "string" },
    { name = "Device", type = "string" },
    { name = "Geo", type = "string" },
    { name = "RemoteAddress", type = "string" },
    { name = "Data", type = "dynamic" },
  ]
}

data "azurerm_client_config" "current" {}

data "azurerm_resource_group" "sentinel" {
  name = var.resource_group_name
}

resource "azuread_application" "descope" {
  display_name = var.application_display_name
}

resource "azuread_service_principal" "descope" {
  client_id = azuread_application.descope.client_id
}

resource "azuread_application_password" "descope" {
  application_id = azuread_application.descope.id

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_log_analytics_workspace_table_custom_log" "descope_audit" {
  name         = local.table_name
  workspace_id = var.log_analytics_workspace_id

  dynamic "column" {
    for_each = local.columns
    content {
      name = column.value.name
      type = column.value.type
    }
  }
}

resource "azapi_resource" "data_collection_rule" {
  type      = "Microsoft.Insights/dataCollectionRules@2024-03-11"
  name      = "dcr-descope-audit-${substr(sha1(var.log_analytics_workspace_id), 0, 13)}"
  parent_id = data.azurerm_resource_group.sentinel.id
  location  = var.location
  body = {
    kind = "Direct"
    properties = {
      streamDeclarations = {
        (local.stream_name) = { columns = [for column in local.columns : { name = column.name, type = lower(column.type) }] }
      }
      destinations = {
        logAnalytics = [{ name = "workspace", workspaceResourceId = var.log_analytics_workspace_id }]
      }
      dataFlows = [{
        streams      = [local.stream_name]
        destinations = ["workspace"]
        transformKql = "source"
        outputStream = local.stream_name
      }]
    }
  }
  response_export_values = ["properties.endpoints.logsIngestion", "properties.immutableId"]
  depends_on             = [azurerm_log_analytics_workspace_table_custom_log.descope_audit]
}

resource "azurerm_role_assignment" "monitoring_metrics_publisher" {
  scope                            = azapi_resource.data_collection_rule.id
  role_definition_id               = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/3913510d-42f4-4e42-8a64-420c390055eb"
  principal_id                     = azuread_service_principal.descope.object_id
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}
```

`outputs.tf`:

```terraform
output "ingestionEndpoint" {
  value = azapi_resource.data_collection_rule.output.properties.endpoints.logsIngestion
}

output "dcrImmutableId" {
  value = azapi_resource.data_collection_rule.output.properties.immutableId
}

output "streamName" {
  value = local.stream_name
}

output "tenantId" {
  value = data.azurerm_client_config.current.tenant_id
}

output "clientId" {
  value = azuread_application.descope.client_id
}

output "clientSecret" {
  value     = azuread_application_password.descope.value
  sensitive = true
}
```

## Usage

Call the module from a configuration where the `azurerm`, `azuread` and `azapi` providers are authenticated against the
subscription and tenant of the workspace, and create the connector in your Descope project from its outputs:

```terraform
provider "azurerm" {
  features {}
}

data "azurerm_log_analytics_workspace" "sentinel" {
  name                = "sentinel-workspace"
  resource_group_name = "sentinel-rg"
}

module "descope_sentinel" {
  source                     = "./descope-sentinel"
  resource_group_name        = data.azurerm_log_analytics_workspace.sentinel.resource_group_name
  log_analytics_workspace_id = data.azurerm_log_analytics_workspace.sentinel.id
  location                   = data.azurerm_log_analytics_workspace.sentinel.location
}

resource "descope_microsoft_sentinel_connector" "sentinel" {
  project_id         = descope_project.myapp.id
  name               = "Microsoft Sentinel"
  ingestion_endpoint = module.descope_sentinel.ingestionEndpoint
  dcr_immutable_id   = module.descope_sentinel.dcrImmutableId
  stream_name        = module.descope_sentinel.streamName
  tenant_id          = module.descope_sentinel.tenantId
  client_id          = module.descope_sentinel.clientId
  client_secret      = module.descope_sentinel.clientSecret
  audit_enabled      = true
}
```

The client secret expires after two years by default. To rotate it before then, run
`terraform apply -replace=module.descope_sentinel.azuread_application_password.descope`. Terraform creates the new
secret and updates the connector with it before it revokes the old secret, so the connector never holds a revoked one.
