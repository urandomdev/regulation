# Repository GuidelinesS

## Project Structure & Module Organization
- `api/` is the Go service: `cmd/regulation` contains the entrypoint, `internal/` holds domain packages, `server/` wires Fiber transport, and `sdk/` plus `sdk/js` surface generated clients.  
- `frontend/` is a Vite + React app (`src/` for UI, `public/` for static assets); install from the repo root so the `@generated/js` workspace resolves.  
- `config/` stores sample JSON configs that are bind-mounted into the API container; use `config.local.json` for developer overrides.  
- `build.sh`, `Dockerfile`, and `docker-compose.yml` orchestrate container builds and a local stack (Postgres, Redis, API).

## Build, Test, and Development Commands
- `npm install` (root) bootstraps all workspaces; follow with targeted commands such as `npm run dev:frontend`.  
- `npm run dev --workspace=frontend` starts the Vite dev server on port 5173; export `VITE_API_URL` if the API host changes.  
- `go run ./cmd/regulation` (inside `api/`) runs the server directly; `docker compose up --build api` is the fastest path to the full stack.  
- `npm run build:frontend` and `npm run build:sdk` emit production bundles; `KO_DOCKER_REPO=<registry> ./build.sh` publishes multi-arch images via `ko`.

## Coding Style & Naming Conventions
- Go sources must stay `gofmt`-ed (tabs, grouped imports). Organize packages by feature (`internal/user`, `internal/virtualaccount`), and keep exported types noun-based (`VirtualAccountService`).  
- TypeScript uses ESLint presets from `eslint.config.js`; keep 2-space indentation, named exports, and PascalCase React components under `src/components`.  
- Use `.ts`/`.tsx` filenames that mirror the default export, e.g., `UserTable.tsx`.

## Testing Guidelines
- Run `go test ./...` inside `api/` before every PR; add focused packages tests in `internal/<feature>/..._test.go`. For coverage, use `go test ./... -coverprofile=coverage.out`.  
- Frontend tests are not scaffolded yet; at minimum run `npm run lint --workspace=frontend` and add Vitest suites under `src/__tests__` when UI logic is added.  
## Commit & Pull Request Guidelines
- Follow the existing Conventional Commit pattern (`feat: add docker-compose`, `fix: logger line`).  
- Each PR should describe the behavior change, list manual test steps, and link the tracking issue. Include screenshots or GIFs for UI-facing edits and reference any schema migrations.  
- Do not mix unrelated fixes; smaller PRs keep the `build.sh` release flow reviewable.

## Security & Configuration Tips
- Never commit secrets; place local overrides in `.env` files or `config/config.local.json` (ignored).  
- The API process reads `CONFIG_DIR` (default `/app/config`)—mirror docker-compose when running locally: `CONFIG_DIR=./config go run ./cmd/regulation`.  
- Rotate database passwords in Postgres by editing `docker-compose.yml` and the corresponding env vars together to avoid drift.
