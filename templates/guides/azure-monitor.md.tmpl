---
page_title: "Azure Monitor Audit Stream - descope Provider"
description: |-
  Create the Azure resources that the Azure Monitor audit stream connector writes to.
---

# Azure Monitor Audit Stream

The Azure Monitor connector streams Descope audit events and troubleshooting logs to a Log Analytics workspace as
OpenTelemetry (OTLP) logs, through the OTLP ingestion endpoint of a data collection rule. The events land in the
`OTelLogs` table of the workspace. This guide provides a Terraform module that creates the Azure side of that setup:

- A data collection endpoint that receives the logs
- A data collection rule that routes the `Microsoft-OTel-Logs` stream into the workspace
- An Entra app registration, service principal and client secret that Descope authenticates as
- A Monitoring Metrics Publisher role assignment for the service principal on the data collection rule

The module uses the `azurerm`, `azuread` and `azapi` providers. The data collection rule is created with `azapi`
because `azurerm_monitor_data_collection_rule` does not support OTLP data sources. OTLP ingestion is an Azure Monitor
preview feature, see [Microsoft's documentation](https://learn.microsoft.com/en-us/azure/azure-monitor/containers/opentelemetry-protocol-ingestion).

## Module

Save the following files in a `descope-azure-monitor` directory.

`variables.tf`:

```terraform
variable "resource_group_name" {
  type        = string
  description = "Resource group to create the data collection endpoint and rule in."
}

variable "log_analytics_workspace_id" {
  type        = string
  description = "Resource ID of the Log Analytics workspace that receives the logs."
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
      version = ">= 4.0"
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
  name_suffix = substr(sha1(var.log_analytics_workspace_id), 0, 13)
}

data "azurerm_client_config" "current" {}

data "azurerm_resource_group" "descope" {
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

resource "azurerm_monitor_data_collection_endpoint" "descope" {
  name                = "dce-descope-${local.name_suffix}"
  resource_group_name = data.azurerm_resource_group.descope.name
  location            = var.location
}

resource "azapi_resource" "data_collection_rule" {
  type      = "Microsoft.Insights/dataCollectionRules@2024-03-11"
  name      = "dcr-descope-${local.name_suffix}"
  parent_id = data.azurerm_resource_group.descope.id
  location  = var.location
  body = {
    properties = {
      dataCollectionEndpointId = azurerm_monitor_data_collection_endpoint.descope.id
      directDataSources = {
        otelLogs = [{
          name                         = "otelLogs"
          streams                      = ["Microsoft-OTel-Logs"]
          enrichWithResourceAttributes = ["*"]
        }]
      }
      destinations = {
        logAnalytics = [{ name = "workspace", workspaceResourceId = var.log_analytics_workspace_id }]
      }
      dataFlows = [{
        streams      = ["Microsoft-OTel-Logs"]
        destinations = ["workspace"]
      }]
    }
  }
  response_export_values = ["properties.immutableId"]
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
output "ingestion_endpoint" {
  value = azurerm_monitor_data_collection_endpoint.descope.logs_ingestion_endpoint
}

output "dcr_immutable_id" {
  value = azapi_resource.data_collection_rule.output.properties.immutableId
}

output "tenant_id" {
  value = data.azurerm_client_config.current.tenant_id
}

output "client_id" {
  value = azuread_application.descope.client_id
}

output "client_secret" {
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

data "azurerm_log_analytics_workspace" "monitor" {
  name                = "monitor-workspace"
  resource_group_name = "monitor-rg"
}

module "descope_azure_monitor" {
  source                     = "./descope-azure-monitor"
  resource_group_name        = data.azurerm_log_analytics_workspace.monitor.resource_group_name
  log_analytics_workspace_id = data.azurerm_log_analytics_workspace.monitor.id
  location                   = data.azurerm_log_analytics_workspace.monitor.location
}

resource "descope_azure_monitor_connector" "monitor" {
  project_id         = descope_project.myapp.id
  name               = "Azure Monitor"
  ingestion_endpoint = module.descope_azure_monitor.ingestion_endpoint
  dcr_immutable_id   = module.descope_azure_monitor.dcr_immutable_id
  tenant_id          = module.descope_azure_monitor.tenant_id
  client_id          = module.descope_azure_monitor.client_id
  client_secret      = module.descope_azure_monitor.client_secret
  audit_enabled      = true

  troubleshoot_log_enabled = true
}
```

The role assignment can take a few minutes to apply, so the first events may be rejected until it does. Query the
events with `OTelLogs | where ServiceName == "descope"`.

The client secret expires after two years by default. To rotate it before then, run
`terraform apply -replace=module.descope_azure_monitor.azuread_application_password.descope`. Terraform creates the new
secret and updates the connector with it before it revokes the old secret, so the connector never holds a revoked one.
