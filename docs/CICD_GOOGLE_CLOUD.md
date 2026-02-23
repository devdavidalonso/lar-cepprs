# CI/CD Google Cloud (Staging e Produção)

## Objetivo
Este projeto agora usa GitHub Actions para:
1. Rodar CI em `push`/`PR` (`main` e `develop`).
2. Fazer deploy do backend no Google Cloud Run, com seleção automática de ambiente:
   - `develop` -> `staging`
   - `main` -> `production`
   - ou manual via `workflow_dispatch`.

Sem depender de `.env` versionado no repositório.

## Workflows criados
- `.github/workflows/ci.yml`
- `.github/workflows/deploy-google-cloud.yml`

## Modelo de configuração segura (sem `.env` no git)
Use **GitHub Environments** (`staging`, `production`) com `Variables` e `Secrets`.

### Variables (por ambiente)
- `GCP_PROJECT_ID`
- `GCP_REGION`
- `GCP_ARTIFACT_REPOSITORY`
- `CLOUD_RUN_BACKEND_SERVICE`
- `BACKEND_ENV_VARS`

Exemplo de `BACKEND_ENV_VARS`:
`APP_ENV=staging,SERVER_PORT=8080,POSTGRES_HOST=10.10.0.3,POSTGRES_PORT=5432,POSTGRES_USER=cecor,POSTGRES_DB=cepr_db,POSTGRES_SSLMODE=disable,SSO_CLIENT_ID=cecor-backend,SSO_REDIRECT_URL=https://api-staging.cecor.hrbsys.tech/*,SSO_AUTH_URL=https://lar-sso-keycloak.hrbsys.tech/realms/cecor/protocol/openid-connect/auth,SSO_TOKEN_URL=https://lar-sso-keycloak.hrbsys.tech/realms/cecor/protocol/openid-connect/token,SSO_USER_INFO_URL=https://lar-sso-keycloak.hrbsys.tech/realms/cecor/protocol/openid-connect/userinfo`

### Secrets (por ambiente)
- `GCP_WORKLOAD_IDENTITY_PROVIDER`
- `GCP_SERVICE_ACCOUNT_EMAIL`
- `BACKEND_SET_SECRETS`

Exemplo de `BACKEND_SET_SECRETS`:
`POSTGRES_PASSWORD=postgres-password:latest,JWT_SECRET=jwt-secret:latest,REFRESH_SECRET=refresh-secret:latest,SSO_CLIENT_SECRET=sso-client-secret:latest,KEYCLOAK_ADMIN_PASSWORD=keycloak-admin-password:latest,SMTP_PASSWORD=smtp-password:latest`

## Pré-requisitos no Google Cloud
1. Criar Artifact Registry (Docker).
2. Criar Cloud Run service (ou deixar workflow criar no primeiro deploy).
3. Criar secrets no Secret Manager.
4. Configurar Workload Identity Federation para GitHub Actions.
5. Dar permissões para service account:
   - `roles/run.admin`
   - `roles/iam.serviceAccountUser`
   - `roles/artifactregistry.writer`
   - `roles/secretmanager.secretAccessor`

## Fluxo operacional simples
1. Desenvolver em branch.
2. Abrir PR -> CI valida backend/frontend.
3. Merge em `develop` -> deploy `staging`.
4. Merge em `main` -> deploy `production`.
5. Se necessário, executar deploy manual em `Actions > Deploy Google Cloud`.

## Observações importantes
1. Valores sensíveis ficam em `Secrets` e Secret Manager, não em `.env` do git.
2. Para testes locais de smoke, use arquivo local `scripts/smoke.env` baseado em `scripts/smoke.env.example`.
3. O frontend já suporta configurações `staging` e `production` via Angular environments.
