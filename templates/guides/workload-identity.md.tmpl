---
page_title: "Workload Identity - descope Provider"
description: |-
  Authenticate the Descope Terraform Provider with a workload's short-lived OIDC token instead of a stored management key.
---

# Workload Identity

A CI/CD workload can authenticate the provider with the short-lived OIDC token its platform issues to it, instead of
a stored management key. There is no secret to store, leak or rotate.

## How It Works

The token is not a credential on its own. Descope authorizes it as one specific **federated management key**, with
exactly the roles that key was given. A federated key holds no secret: it trusts a single OIDC issuer, and only
accepts tokens from that issuer whose subject matches the patterns it was configured with.

The token names the key it acts as in its `aud` claim, so the audience must be set when the token is requested.

## Create a Federated Management Key

In [Company Settings](https://app.descope.com/settings/company), create a management key with the
**Authentication type** set to **Federated - Workload Identity**, and configure:

- **Trusted Issuer** - The HTTPS issuer URL of the platform that issues the workload's tokens. Its signing keys are
  discovered from `<issuer>/.well-known/openid-configuration`, so the issuer must be publicly reachable.
- **Allowed Subjects** - Regular expressions matched against the whole `sub` claim. Only tokens whose subject matches
  one of them are accepted. Patterns that would match any subject are rejected.
- **Max Token Timeout** - The longest token lifetime the key accepts, between 1 and 15 minutes. Tokens issued with a
  longer lifetime are rejected.
- **Roles** - The same roles as any other management key, for example project roles on the projects your
  configuration manages.

Once the key is created, copy its **Audience (aud)** value.

## Configure the Provider

Request a token from the workload's platform with the audience you copied, exactly as shown in the console, and pass
it to the provider in place of a management key. No other configuration changes:

```shell
export DESCOPE_MANAGEMENT_KEY="<workload identity token>"
terraform plan
```

The `management_key` provider attribute accepts the token as well.

## Token Lifetime

The provider sends the token with every request it makes, and each request is checked against the token's expiry.
Request a fresh token for each Terraform run, and keep in mind that a run which outlives its token fails partway
through with an authentication error. Run it again with a new token to continue.

## Troubleshooting

If requests are rejected with an authentication error, check the token's claims against the key's configuration:

- `iss` is exactly the key's **Trusted Issuer**
- `aud` contains the key's **Audience (aud)** value
- `sub` matches one of the key's **Allowed Subjects**
- The token's lifetime, from `iat` to `exp`, is within the key's **Max Token Timeout**, and it has not expired
