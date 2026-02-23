#!/usr/bin/env bash
# Smoke test RBAC com Keycloak como fonte única de token/roles.

set -u

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="${SMOKE_ENV_FILE:-${SCRIPT_DIR}/smoke.env}"
if [[ -f "${ENV_FILE}" ]]; then
  set -a
  # shellcheck disable=SC1090
  source "${ENV_FILE}"
  set +a
fi

PROFILE="${SMOKE_PROFILE:-staging}"
if [[ "${PROFILE}" == "production" ]]; then
  DEFAULT_BASE_URL="${PRODUCTION_BASE_URL:-http://localhost:8082}"
  DEFAULT_TOKEN_URL="${PRODUCTION_KEYCLOAK_TOKEN_URL:-http://localhost:8081/realms/cepprs/protocol/openid-connect/token}"
  DEFAULT_CLIENT_ID="${PRODUCTION_KEYCLOAK_CLIENT_ID:-lar-cepprs-backend}"
  DEFAULT_CLIENT_SECRET="${PRODUCTION_KEYCLOAK_CLIENT_SECRET:-}"
  DEFAULT_ADMIN_USER="${PRODUCTION_ADMIN_USER:-admin.cecor}"
  DEFAULT_ADMIN_PASS="${PRODUCTION_ADMIN_PASS:-admin123}"
  DEFAULT_PROF_USER="${PRODUCTION_PROF_USER:-prof.maria}"
  DEFAULT_PROF_PASS="${PRODUCTION_PROF_PASS:-prof123}"
  DEFAULT_STUDENT_USER="${PRODUCTION_STUDENT_USER:-aluno.pedro}"
  DEFAULT_STUDENT_PASS="${PRODUCTION_STUDENT_PASS:-aluno123}"
else
  DEFAULT_BASE_URL="${STAGING_BASE_URL:-http://localhost:8082}"
  DEFAULT_TOKEN_URL="${STAGING_KEYCLOAK_TOKEN_URL:-http://localhost:8081/realms/cepprs/protocol/openid-connect/token}"
  DEFAULT_CLIENT_ID="${STAGING_KEYCLOAK_CLIENT_ID:-lar-cepprs-backend}"
  DEFAULT_CLIENT_SECRET="${STAGING_KEYCLOAK_CLIENT_SECRET:-}"
  DEFAULT_ADMIN_USER="${STAGING_ADMIN_USER:-admin.cecor}"
  DEFAULT_ADMIN_PASS="${STAGING_ADMIN_PASS:-admin123}"
  DEFAULT_PROF_USER="${STAGING_PROF_USER:-prof.maria}"
  DEFAULT_PROF_PASS="${STAGING_PROF_PASS:-prof123}"
  DEFAULT_STUDENT_USER="${STAGING_STUDENT_USER:-aluno.pedro}"
  DEFAULT_STUDENT_PASS="${STAGING_STUDENT_PASS:-aluno123}"
fi

BASE_URL="${BASE_URL:-${DEFAULT_BASE_URL}}"
API_BASE="${BASE_URL}/api/v1"

KEYCLOAK_TOKEN_URL="${KEYCLOAK_TOKEN_URL:-${DEFAULT_TOKEN_URL}}"
KEYCLOAK_CLIENT_ID="${KEYCLOAK_CLIENT_ID:-${DEFAULT_CLIENT_ID}}"
KEYCLOAK_CLIENT_SECRET="${KEYCLOAK_CLIENT_SECRET:-${DEFAULT_CLIENT_SECRET}}"

ADMIN_USER="${ADMIN_USER:-${DEFAULT_ADMIN_USER}}"
ADMIN_PASS="${ADMIN_PASS:-${DEFAULT_ADMIN_PASS}}"
PROF_USER="${PROF_USER:-${DEFAULT_PROF_USER}}"
PROF_PASS="${PROF_PASS:-${DEFAULT_PROF_PASS}}"
STUDENT_USER="${STUDENT_USER:-${DEFAULT_STUDENT_USER}}"
STUDENT_PASS="${STUDENT_PASS:-${DEFAULT_STUDENT_PASS}}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASSED=0
FAILED=0

print_header() {
  echo "======================================"
  echo "🧪 CECOR Smoke RBAC (Keycloak)"
  echo "======================================"
  echo "API:       ${API_BASE}"
  echo "Token URL: ${KEYCLOAK_TOKEN_URL}"
  echo "Client ID: ${KEYCLOAK_CLIENT_ID}"
  echo "Profile:   ${PROFILE}"
  echo ""
}

extract_access_token() {
  sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p'
}

get_token() {
  local username="$1"
  local password="$2"
  local body
  local token

  if [[ -n "${KEYCLOAK_CLIENT_SECRET}" ]]; then
    body="$(curl -sS -X POST "${KEYCLOAK_TOKEN_URL}" \
      -H "Content-Type: application/x-www-form-urlencoded" \
      --data-urlencode "grant_type=password" \
      --data-urlencode "client_id=${KEYCLOAK_CLIENT_ID}" \
      --data-urlencode "client_secret=${KEYCLOAK_CLIENT_SECRET}" \
      --data-urlencode "username=${username}" \
      --data-urlencode "password=${password}" || true)"
  else
    body="$(curl -sS -X POST "${KEYCLOAK_TOKEN_URL}" \
      -H "Content-Type: application/x-www-form-urlencoded" \
      --data-urlencode "grant_type=password" \
      --data-urlencode "client_id=${KEYCLOAK_CLIENT_ID}" \
      --data-urlencode "username=${username}" \
      --data-urlencode "password=${password}" || true)"
  fi

  token="$(printf '%s' "${body}" | extract_access_token)"
  printf '%s' "${token}"
}

http_code() {
  local method="$1"
  local path="$2"
  local token="$3"

  if [[ -n "${token}" ]]; then
    curl -s -o /dev/null -w "%{http_code}" -X "${method}" "${API_BASE}${path}" \
      -H "Authorization: Bearer ${token}" 2>/dev/null || true
  else
    curl -s -o /dev/null -w "%{http_code}" -X "${method}" "${API_BASE}${path}" 2>/dev/null || true
  fi
}

assert_code() {
  local method="$1"
  local path="$2"
  local token="$3"
  local expected="$4"
  local label="$5"
  local code

  code="$(http_code "${method}" "${path}" "${token}")"
  if [[ "${code}" == "${expected}" ]]; then
    echo -e "${GREEN}✓${NC} ${label} (HTTP ${code})"
    PASSED=$((PASSED + 1))
  else
    echo -e "${RED}✗${NC} ${label} (HTTP ${code}, esperado ${expected})"
    FAILED=$((FAILED + 1))
  fi
}

assert_code_in() {
  local method="$1"
  local path="$2"
  local token="$3"
  local expected_csv="$4"
  local label="$5"
  local code

  code="$(http_code "${method}" "${path}" "${token}")"
  IFS=',' read -r -a allowed <<< "${expected_csv}"
  for value in "${allowed[@]}"; do
    if [[ "${code}" == "${value}" ]]; then
      echo -e "${GREEN}✓${NC} ${label} (HTTP ${code})"
      PASSED=$((PASSED + 1))
      return
    fi
  done

  echo -e "${RED}✗${NC} ${label} (HTTP ${code}, esperado um de ${expected_csv})"
  FAILED=$((FAILED + 1))
}

print_summary() {
  echo ""
  echo "======================================"
  echo "📊 RESULTADO SMOKE RBAC"
  echo "======================================"
  echo -e "${GREEN}Passaram: ${PASSED}${NC}"
  echo -e "${RED}Falharam: ${FAILED}${NC}"
  echo ""

  if [[ "${FAILED}" -eq 0 ]]; then
    echo -e "${GREEN}Smoke RBAC aprovado.${NC}"
    exit 0
  fi

  echo -e "${YELLOW}Smoke RBAC com falhas.${NC}"
  exit 1
}

print_header

echo "Obtendo tokens de teste..."
ADMIN_TOKEN="$(get_token "${ADMIN_USER}" "${ADMIN_PASS}")"
PROF_TOKEN="$(get_token "${PROF_USER}" "${PROF_PASS}")"
STUDENT_TOKEN="$(get_token "${STUDENT_USER}" "${STUDENT_PASS}")"

if [[ -z "${ADMIN_TOKEN}" ]]; then
  echo -e "${RED}Falha ao obter token admin (${ADMIN_USER}).${NC}"
  echo "Verifique usuário/senha, client e se Direct Access Grants está habilitado no client."
  exit 2
fi
if [[ -z "${PROF_TOKEN}" ]]; then
  echo -e "${RED}Falha ao obter token professor (${PROF_USER}).${NC}"
  exit 2
fi
if [[ -z "${STUDENT_TOKEN}" ]]; then
  echo -e "${RED}Falha ao obter token aluno (${STUDENT_USER}).${NC}"
  exit 2
fi

echo -e "${GREEN}Tokens obtidos com sucesso.${NC}"
echo ""

# Auth baseline
assert_code "GET" "/auth/verify" "" "401" "Sem token deve negar /auth/verify"
assert_code "GET" "/auth/verify" "${ADMIN_TOKEN}" "200" "Admin autenticado em /auth/verify"
assert_code "GET" "/auth/verify" "${PROF_TOKEN}" "200" "Professor autenticado em /auth/verify"
assert_code "GET" "/auth/verify" "${STUDENT_TOKEN}" "200" "Aluno autenticado em /auth/verify"

# RBAC admin route
assert_code "GET" "/users/teachers" "${ADMIN_TOKEN}" "200" "Admin pode acessar /users/teachers"
assert_code "GET" "/users/teachers" "${PROF_TOKEN}" "403" "Professor não pode acessar /users/teachers"
assert_code "GET" "/users/teachers" "${STUDENT_TOKEN}" "403" "Aluno não pode acessar /users/teachers"

# Role-specific portals (aceita 404 quando perfil ainda não provisionado no banco local)
assert_code_in "GET" "/teacher/dashboard" "${PROF_TOKEN}" "200,404" "Professor acessa portal professor"
assert_code_in "GET" "/student/dashboard" "${STUDENT_TOKEN}" "200,404" "Aluno acessa portal aluno"

print_summary
