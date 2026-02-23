-- 14-create-multi-program-structure.sql
-- Estrutura multi-programa para Centro Educacional:
-- - Centro educacional
-- - Programas (Semear, Voar, Cecor)
-- - Vinculo aluno <-> programa
-- - Referencia opcional de programa em cursos e localizacoes

CREATE TABLE IF NOT EXISTS public.educational_centers (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    code TEXT NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.programs (
    id BIGSERIAL PRIMARY KEY,
    center_id BIGINT NOT NULL REFERENCES public.educational_centers(id),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.student_programs (
    id BIGSERIAL PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES public.students(id),
    program_id BIGINT NOT NULL REFERENCES public.programs(id),
    status TEXT NOT NULL DEFAULT 'active',
    entry_date TIMESTAMPTZ NULL,
    exit_date TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT student_program_unique UNIQUE (student_id, program_id)
);

ALTER TABLE public.courses
    ADD COLUMN IF NOT EXISTS program_id BIGINT REFERENCES public.programs(id);

ALTER TABLE public.locations
    ADD COLUMN IF NOT EXISTS center_id BIGINT REFERENCES public.educational_centers(id),
    ADD COLUMN IF NOT EXISTS program_id BIGINT REFERENCES public.programs(id);

-- Seed idempotente dos valores iniciais
INSERT INTO public.educational_centers (name, code, is_active, created_at, updated_at)
VALUES ('Centro Educacional Prof. Paulo Rossi Severino', 'CEPROS', TRUE, NOW(), NOW())
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
    updated_at = NOW();

INSERT INTO public.programs (center_id, code, name, is_active, created_at, updated_at)
SELECT ec.id, p.code, p.name, TRUE, NOW(), NOW()
FROM public.educational_centers ec
JOIN (
    VALUES
        ('SEMEAR', 'Semear'),
        ('VOAR', 'Voar'),
        ('CECOR', 'Cecor')
) AS p(code, name) ON TRUE
WHERE ec.code = 'CEPROS'
ON CONFLICT (code) DO UPDATE
SET center_id = EXCLUDED.center_id,
    name = EXCLUDED.name,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

