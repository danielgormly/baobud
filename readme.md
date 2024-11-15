# Baobud

[![Release Baobud](https://github.com/danielgormly/baobud/actions/workflows/build.yml/badge.svg)](https://github.com/danielgormly/baobud/actions/workflows/build.yml)

Generate [OpenBao](https://openbao.org/)/[Vault](https://www.hashicorp.com/products/vault) policies from [Consul Template](https://github.com/hashicorp/consul-template) templates. Baobud dynamically evaluates the template to determine all Vault requests.

## Limitations & Caveats
- Baobud does not support evaluating Consul nor Nomad requests.
- This is currently built with the Consul Template SDK, that uses the Vault SDK. If OpenBao/Vault API diverges, this will break for OpenBao.
- Incorrect auth will give you an ugly runtime error.

## Usage example

**input: template.toml.tmpl**
```
{{ with secret "secret/prod/example" }}
  EXAMPLE_SECRET: {{ .Data.EXAMPLE_SECRET }}
  EXAMPLE_CONFIG: {{ .Data.EXAMPLE_CONFIG }}
{{ end }}
```

### Usage
```bash
# Create policy & write to stdout
baobud -f template.ctmpl
# path "secret/prod/example" {
#  capabilities = ["read"]
# }

# Create policy & write to file
baobud -f template.ctmpl -o policy.hcl
```

## Other commands
```bash
baobud version # prints version
baobud help # prints help info
```

## Installation
```bash
# Linux
curl -Lo baobud https://github.com/danielgormly/baobud/releases/download/v0.0.1-alpha-11/baobud-linux-amd64
# MacOS (ARM-based)
curl -Lo baobud https://github.com/danielgormly/baobud/releases/download/v0.0.1-alpha-11/baobud-darwin-arm64
# Install
chmod +x baobud
mv baobud /usr/local/bin/
baobud version
```
