-- 15-create-teacher-programs.sql
-- Normalizacao do vinculo professor <-> programa.

CREATE TABLE IF NOT EXISTS public.teacher_programs (
    id BIGSERIAL PRIMARY KEY,
    teacher_id BIGINT NOT NULL REFERENCES public.teachers(id),
    program_id BIGINT NOT NULL REFERENCES public.programs(id),
    role TEXT NOT NULL DEFAULT 'teacher',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT teacher_program_unique UNIQUE (teacher_id, program_id)
);

-- Backfill com base em cursos já atribuídos aos professores
INSERT INTO public.teacher_programs (teacher_id, program_id, role, is_active, created_at, updated_at)
SELECT DISTINCT tc.teacher_id, c.program_id, 'teacher', TRUE, NOW(), NOW()
FROM public.teacher_courses tc
JOIN public.courses c ON c.id = tc.course_id
WHERE c.program_id IS NOT NULL
ON CONFLICT (teacher_id, program_id) DO NOTHING;

