.PHONY: dev
dev:
	cd ui && npm run dev

.PHONY: generate
generate:
	go run -mod=mod ./ent/entc.go && go generate ./... && go mod tidy && cd ui && pnpm run codegen

.PHONY: test
test:
	go test -v -cover ./...
