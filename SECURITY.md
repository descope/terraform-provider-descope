# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest  | Yes       |

## Reporting a Vulnerability

If you discover a security vulnerability in this fork, please report it responsibly.

**Do not open a public GitHub issue for security vulnerabilities.**

Instead, please use [GitHub Private Vulnerability Reporting](https://github.com/jamescrowley321/terraform-provider-descope/security/advisories/new) to submit your report.

Include:
- A description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

You can expect an initial response within 72 hours. We will work with you to understand and address the issue before any public disclosure.

## Security Scanning

This repository uses the following security tools:

| Tool | Purpose |
|------|---------|
| [CodeQL](https://codeql.github.com/) | Static analysis (SAST) |
| [gosec](https://github.com/securego/gosec) | Go-specific security analysis |
| [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) | Go vulnerability database scanning |
| [Dependabot](https://docs.github.com/en/code-security/dependabot) | Automated dependency updates |
| [gitleaks](https://github.com/gitleaks/gitleaks) | Secret detection |
| [GitHub Secret Scanning](https://docs.github.com/en/code-security/secret-scanning) | Repository-level secret detection |
| [OpenSSF Scorecard](https://securityscorecards.dev/) | Supply chain security assessment |

## Vulnerability Response SLA

| Severity | Response Target |
|----------|-----------------|
| Critical | 24 hours |
| High     | 7 days |
| Medium   | 30 days |
| Low      | Next release cycle |

## Upstream Vulnerabilities

For vulnerabilities in the upstream [descope/terraform-provider-descope](https://github.com/descope/terraform-provider-descope) project, please report them directly to the upstream maintainers. We will incorporate upstream security fixes as part of our regular sync process.
