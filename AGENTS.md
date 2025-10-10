# Repository Guidelines

## Project Structure & Module Organization
The Go API entry point lives in `cmd/kubernetes-auditing-dashboard`, with shared services under `pkg/services`. GraphQL schemas, resolvers, and tests reside in `gql`, while Ent entity definitions sit in `ent/schema` and generate code into `ent`. The Next.js client is under `ui` (`ui/pages`, `ui/modules`, `ui/styles`, `ui/public`). Kubernetes bootstrap assets for local clusters are in `script/kube-apiserver-config`.

## Build, Test, and Development Commands
- `make dev` (root): runs the Next.js dev server, GraphQL codegen watch, and Go backend for a full-stack loop.
- `make generate`: regenerates Ent code, executes `go generate ./...`, tidies modules, and refreshes typed GraphQL documents.
- `cd ui && npm run dev`: mirrors `make dev` when you only need the front-end workflow.
- `go test ./...`: executes Go unit tests, including pagination coverage in `gql/*_test.go`.
- `cd ui && npm run lint`: applies the bundled Next.js/ESLint checks.

## Coding Style & Naming Conventions
Run `gofmt` (or `go fmt ./...`) before committing backend changes; exported Go symbols stay PascalCase, internal helpers camelCase, and avoid package-level underscores. Keep Ent schemas declarative and colocated in `ent/schema/<entity>.go`. In the UI, rely on the bundled ESLint config; React components and hooks follow `PascalCase` and `useCamelCase` respectively, and co-locate shared UI logic under `ui/modules`. Favor Tailwind utility classes already used in `ui/styles`.

## Testing Guidelines
Add Go tests alongside the package under test (e.g., `pkg/<area>/*_test.go`). Extend GraphQL coverage by mirroring the pagination test pattern in `gql/pagination_test.go` whenever you add resolvers. For UI tweaks, pair visual changes with screenshots in PRs and run `npm run lint`; introduce Testing Library or Playwright specs under `ui/__tests__` when creating interactive flows.

## Commit & Pull Request Guidelines
Follow the Conventional Commits pattern visible in the log (`feat:`, `fix:`, `chore:`). Keep the subject imperative and under 72 characters; details belong in the body. PRs should describe the change, reference related issues, list manual or automated checks (e.g., `go test ./...`, `npm run lint`), and attach UI screenshots when the front end shifts. Request review once CI or local checks pass.

## AI Assistant Workflow
When using AI assistants (like Claude Code) for development:
- **NEVER run `git commit`** - Always stage changes with `git add` and let the human developer review and commit
- Show changes with `git status` and `git diff` when work is complete
- Run build/lint checks (`npm run build`, `npm run lint`, `go test ./...`) to verify changes compile
- The developer will review all staged changes before committing

## Security & Configuration Tips
Treat `script/kube-apiserver-config` as the source of truth for Minikube auditing setup—update both the manifest patch and README instructions together. Never commit real auditing data; refresh or redact `data.db` before pushing. Store sensitive connection strings in environment variables or Kubernetes secrets rather than hardcoding them in Go or Next.js modules.
