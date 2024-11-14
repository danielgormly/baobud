## Development & testing rough notes
```bash
bao server -dev -dev-root-token-id=dev
export BAO_ADDR='http://127.0.0.1:8200'
export BAO_TOKEN='dev'
bao secrets enable -path=kv2 -version=2 kv
bao kv put kv2/test foo="bar" bar="foo"
bao kv get kv2/test
consul-template -template="template.ctmpl" -vault-addr="http://127.0.0.1:8200" \
    -vault-token="dev" -vault-renew-token=false -once -dry
```

## Release
git tag v0.0.1
git push origin v0.0.1
