```bash
oapi-codegen -generate "types,client" \
  -package osv \
  -o internal/api/osv/client.gen.go \
  docs/openapi/osv/osv_service_v1.yaml
```