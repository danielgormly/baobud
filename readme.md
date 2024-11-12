# Baobud

Generate [OpenBao](https://openbao.org/)/[Vault](https://www.hashicorp.com/products/vault) policies from [Consul Template](https://github.com/hashicorp/consul-template) templates. Baobud is not a static tool so will

## Limitations & Caveats
- Baobud does not support Consul.
- This is currently built with the Consul Template SDK, that uses the Vault SDK. If OpenBao/Vault API diverges, this will break for OpenBao.

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
> path "secret/prod/example" {
>   capabilities = ["read"]
> }

# Create policy & write to file
baobud -f template.ctmpl -o policy.hcl
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
```bash
# Nested dynamic key test!
bao server -dev -dev-root-token-id=dev
export BAO_ADDR='http://127.0.0.1:8200'
export BAO_TOKEN='dev'
bao secrets enable -path=kv2 -version=2 kv

# Nested example
bao kv put kv2/app name="myapp" environment="prod"
bao kv put kv2/env/prod type="production" config="debug=false"
consul-template -template="nested.ctmpl" -vault-addr="http://127.0.0.1:8200" -vault-token="dev" -vault-renew-token=false -once -dry
baobud kv

export BAOBUD_DEBUG=1 # see debug info
go run ./cmd -f ./test/nested.ctmpl # compile template
make build # build

## Dynamic test (TODO: Encapsulate)
TEST_DYNAMIC=yo go run ./cmd -f ./test/dynamic.ctmpl

## Nested
consul-template -template="ultra.ctmpl" -vault-addr="http://127.0.0.1:8200" -vault-token="s.JQLssCtAbBiJkDe0q2HyL6v0" -vault-renew-token=false -once -dr
```
