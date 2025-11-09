# 設計書

# 概要

# 目的と課題

# 外部的に必要な機能

- [ ] 

# 内部的に必要な機能

## internal/github/api

- [ ]  リリース一覧取得機能
- [ ]  期間を指定し、特定の期間にマージされたPR一覧取得機能
- [ ] バージョン情報を基に、脆弱性の

## internal/github/osv

- [ ] 脆弱性の一覧取得API
- [ ] CVE番号から、脆弱性が導入されたコミット、修正コミットの予測を取得する
- [ ] CVE番号から脆弱性が含まれているバージョン情報を取得する




```bash
oapi-codegen -generate "types,client" \
  -package osv \
  -o internal/api/osv/client.gen.go \
  docs/openapi/osv/osv_service_v1.yaml
```