# Status de Progresso - 22/02/2026

## Resumo executivo

Status atual da validacao: `GO para continuidade da V1 (ambiente local/homologacao)`.

Foram concluidos os pontos criticos que bloqueavam testes manuais:
1. RBAC por perfil (admin/professor/aluno) com navegacao separada.
2. Fluxo de login Keycloak consistente no frontend.
3. Rotas protegidas negando acesso indevido (`/students` para professor/aluno).
4. `students` e `teachers` respondendo com `200` no frontend.
5. Seed local idempotente para popular professores e alunos.
6. Correcao de schema legado no Postgres (`students.special_needs` e compatibilidade).
7. Correcao do proxy frontend para backend local.
8. Correcao do erro de asset `icon-144x144.png`.

## Evidencias de validacao

1. Quick preflight da API aprovado (`4/4`).
2. Smoke RBAC aprovado (`9/9`).
3. Teste manual:
   - Admin lista alunos com dados.
   - Admin lista professores com dados.
   - Professor sem acesso a area administrativa.
   - Aluno com portal proprio.
4. Requests no browser para `students` e `teachers` com status `200`.
5. Fluxo de homologacao funcional (admin):
   - professor criado com `201`.
   - curso criado com `201`.
   - aluno criado com `201`.
   - matricula criada com `201`.
6. Regra de negocio de idade validada:
   - matricula bloqueada para aluno < 12 anos.
   - apos ajuste de data de nascimento, matricula permitida.
7. Portal do professor:
   - turma e alunos aparecem apos vinculacao em `teacher_courses`.
   - tela de perfil do aluno em rota de professor ainda em placeholder ("Fase 5").

## Comandos oficiais usados na homologacao local

```bash
cd /home/david-alonso/Projetos/lar-sso && docker compose up -d
cd /home/david-alonso/Projetos/cecor && make dev-infra-up
cd /home/david-alonso/Projetos/cecor && make dev-backend
cd /home/david-alonso/Projetos/cecor && make dev-frontend
cd /home/david-alonso/Projetos/cecor && make dev-seed
cd /home/david-alonso/Projetos/cecor && ./scripts/quick_api_test.sh
cd /home/david-alonso/Projetos/cecor && ./scripts/smoke_rbac_keycloak.sh
```

## Mudancas tecnicas relevantes consolidadas

1. Frontend:
   - Guards/rotas por perfil.
   - Sidebar por perfil sem duplicidade de menu.
   - Redirect por papel apos login.
   - Proxy `/api` apontando para backend local (`8082`).
   - Fallback e icones PWA corrigidos.
2. Backend:
   - Compatibilidade de schema para tabela `students`.
   - Ajustes em rotas e estabilidade de handlers.
   - Seed local para ambiente de desenvolvimento.
3. DevEx:
   - Modo hibrido: app local + infra em Docker.
   - Targets de `make` para operacao diaria.

## Pendencias restantes (nao bloqueantes para ciclo atual)

1. Consolidar dados de seed com maior massa (cursos/matriculas/presencas).
2. Expandir suite automatizada de regressao frontend.
3. Fechar checklist de release e versionamento final da V1.
4. Completar UX/fluxo do professor para lancamento de presenca direto na tela de turma/alunos.
5. Completar menu do professor (links e atalhos operacionais para chamada).
6. Corrigir calculo de carga horaria do curso (workload negativo observado em homologacao).

## Decisoes e implementacoes consolidadas (memoria de produto)

1. O projeto passa a adotar o nome de produto `LAR CEPPRS` (renomeacao progressiva de branding e deploy IDs).
2. O menu admin deve expor cadastro de dominio para `Locais e Salas`.
3. A capacidade de matriculas deve refletir a capacidade operacional:
   - regra aplicada no backend para bloquear matricula quando `max_students` do curso atingir limite.
4. Staging e production devem usar variaveis/secrets em GitHub Environments (sem `.env` versionado).
5. Validacao critica da entidade de dominio `locations` documentada em:
   - `docs/LOCATIONS_ENTIDADE_VALIDACAO.md`
