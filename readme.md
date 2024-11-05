# Baobud

Generate [OpenBao](https://openbao.org/)/[Vault](https://www.hashicorp.com/products/vault) policies from [Consul Template](https://github.com/hashicorp/consul-template) templates.

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
baobud template.ctmpl
> path "secret/prod/example" {
>   capabilities = ["read"]
> }

# Create policy & write to file
baobud template.ctmpl -o policy.hcl
```

## Other commands
```bash
baobud version # prints version
baobud help # prints help info
```

## Installation (MacOS ARM)
```
curl baobud
chmod +x baobud
mv baobud /usr/bin/local
baobud version
```

## Installation (Linux)
```
curl baobud
chmod +x baobud
mv baobud /usr/bin
baobud version
```

## Development
`go run main.go version`
