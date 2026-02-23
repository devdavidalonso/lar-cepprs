# Entidade Locations - Validacao e Realidade Atual

**Data:** 22/02/2026  
**Projeto:** LAR CEPPRS

## 1. O que a entidade `locations` deve representar

`locations` deve representar o dominio fisico de oferta:
1. Sala/laboratorio/espaco.
2. Capacidade real do espaco.
3. Recursos disponiveis (projetor, computadores, etc).
4. Disponibilidade para uso.

Objetivo de negocio: a capacidade da sala deve limitar o numero de alunos matriculados.

## 2. Estado atual encontrado (AS-IS)

### 2.1 Backend (modelo)

Arquivo: `backend/internal/models/location.go`

A tabela/modelo existe com:
1. `name`
2. `capacity`
3. `resources`
4. `isActive`

### 2.2 Frontend (telas/rotas)

Arquivos:
1. `frontend/src/app/features/locations/locations.routes.ts`
2. `frontend/src/app/features/locations/components/location-form.component.ts`
3. `frontend/src/app/core/services/location.service.ts`

Observacoes:
1. Existe tela de cadastro de local/sala.
2. O service de locations atualmente usa **dados mockados em memoria**.
3. A tela salva localmente (simulacao), sem persistencia real em API.

### 2.3 Integracao com Cursos

Arquivo: `frontend/src/app/features/courses/components/course-form.component.ts`

Observacoes:
1. O form de curso pede `locationId` e exibe capacidade da sala no select.
2. Porem o backend de `courses` ainda nao persiste `locationId` no modelo `Course`.
3. Resultado: o vinculo curso-sala ainda nao e fonte unica de verdade.

### 2.4 Integracao com Matriculas

Arquivo: `backend/internal/service/enrollments/service.go`

Observacoes:
1. Ja existe bloqueio por capacidade de curso (`maxStudents`).
2. Ainda falta forcar que `maxStudents` respeite capacidade da sala vinculada.

## 3. Diagnostico critico

Hoje a entidade `locations` existe, mas ainda nao governa o fluxo fim-a-fim.

Gaps:
1. Sem CRUD backend oficial para locations no fluxo principal.
2. Sem persistencia real de locations no frontend (mock).
3. Sem vinculo persistido e validado entre `course` e `location`.
4. Regra de capacidade ainda ancorada em `course.maxStudents`, nao na sala.

## 4. Regra de negocio alvo (TO-BE)

Regra recomendada:
1. Curso deve ter sala vinculada (`locationId`) quando for presencial.
2. `course.maxStudents` nao pode ser maior que `location.capacity`.
3. Na matricula, limite efetivo = `min(course.maxStudents, location.capacity)`.

## 5. Criterios de aceite da entidade (homologacao)

1. Admin cadastra sala com capacidade e recursos.
2. Sala aparece no cadastro de curso via API real (nao mock).
3. Curso nao salva se `maxStudents > capacity da sala`.
4. Matricula bloqueia quando atingir limite da sala/curso.
5. Professor enxerga local correto na turma/sessao.

## 6. Proximo passo tecnico recomendado

1. Implementar CRUD backend de `locations` (`/api/v1/locations`).
2. Migrar frontend `LocationService` de mock para HTTP.
3. Persistir `locationId` em `Course` (ou consolidar via `CourseClass` como fonte unica).
4. Adicionar validacao backend no create/update de curso para capacidade da sala.
