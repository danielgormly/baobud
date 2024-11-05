# Baobud

Generate [OpenBao](https://openbao.org/)/[Vault](https://www.hashicorp.com/products/vault) policies from [Consul Template](https://github.com/hashicorp/consul-template) templates.

## Usage example

*input: template.toml.tmpl*
```
{{ with secret "secret/prod/example" }}
  EXAMPLE_SECRET: {{ .Data.EXAMPLE_SECRET }}
  EXAMPLE_CONFIG: {{ .Data.EXAMPLE_CONFIG }}
{{ end }}
```

```bash
baobud template.ctmpl > policy.hcl
#> path "secret/data/prod/example" {
#>   capabilities = ["read"]
#> }
# All secret engines assume v2 secret engine unless v1 path(s) are matched. Note you can use multiple --v1-prefix arguments
baobud --v1-prefix="secret" template.ctmpl > policy.hcl
#> path "secret/prod/example" {
#>   capabilities = ["read"]
#> }
```

## Other commands
```bash
baobud version # prints version
baobud help # prints help info
```

## Installation (MacOS)
```
```

## Installation (Linux)
```
```

## Development
`go run main.go version`
