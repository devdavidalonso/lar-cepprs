-- 17-normalize-multi-program-constraints.sql
-- Normaliza nomes e existência de constraints/índices do contexto multi-programa.
-- Idempotente e segura para ambientes com histórico legado.

BEGIN;

DO $$
BEGIN
    -- educational_centers (legacy gorm names -> canonical)
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'uni_educational_centers_name'
          AND conrelid = 'public.educational_centers'::regclass
    ) AND NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'educational_centers_name_key'
          AND conrelid = 'public.educational_centers'::regclass
    ) THEN
        ALTER TABLE public.educational_centers
            RENAME CONSTRAINT uni_educational_centers_name TO educational_centers_name_key;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'uni_educational_centers_code'
          AND conrelid = 'public.educational_centers'::regclass
    ) AND NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'educational_centers_code_key'
          AND conrelid = 'public.educational_centers'::regclass
    ) THEN
        ALTER TABLE public.educational_centers
            RENAME CONSTRAINT uni_educational_centers_code TO educational_centers_code_key;
    END IF;

    -- programs (legacy gorm names -> canonical)
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'uni_programs_code'
          AND conrelid = 'public.programs'::regclass
    ) AND NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'programs_code_key'
          AND conrelid = 'public.programs'::regclass
    ) THEN
        ALTER TABLE public.programs
            RENAME CONSTRAINT uni_programs_code TO programs_code_key;
    END IF;

    -- student_programs / teacher_programs (legacy unique index names -> canonical)
    IF EXISTS (
        SELECT 1
        FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE n.nspname = 'public'
          AND c.relkind = 'i'
          AND c.relname = 'idx_student_program_unique'
    ) AND NOT EXISTS (
        SELECT 1
        FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE n.nspname = 'public'
          AND c.relkind = 'i'
          AND c.relname = 'student_program_unique'
    ) THEN
        ALTER INDEX public.idx_student_program_unique RENAME TO student_program_unique;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE n.nspname = 'public'
          AND c.relkind = 'i'
          AND c.relname = 'idx_teacher_program_unique'
    ) AND NOT EXISTS (
        SELECT 1
        FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE n.nspname = 'public'
          AND c.relkind = 'i'
          AND c.relname = 'teacher_program_unique'
    ) THEN
        ALTER INDEX public.idx_teacher_program_unique RENAME TO teacher_program_unique;
    END IF;

    -- Guarantees: expected unique constraints exist (for fresh/legacy bases).
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'educational_centers_name_key'
          AND conrelid = 'public.educational_centers'::regclass
          AND contype = 'u'
    ) THEN
        ALTER TABLE public.educational_centers
            ADD CONSTRAINT educational_centers_name_key UNIQUE (name);
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'educational_centers_code_key'
          AND conrelid = 'public.educational_centers'::regclass
          AND contype = 'u'
    ) THEN
        ALTER TABLE public.educational_centers
            ADD CONSTRAINT educational_centers_code_key UNIQUE (code);
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'programs_code_key'
          AND conrelid = 'public.programs'::regclass
          AND contype = 'u'
    ) THEN
        ALTER TABLE public.programs
            ADD CONSTRAINT programs_code_key UNIQUE (code);
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'student_program_unique'
          AND conrelid = 'public.student_programs'::regclass
          AND contype = 'u'
    ) THEN
        ALTER TABLE public.student_programs
            ADD CONSTRAINT student_program_unique UNIQUE (student_id, program_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'teacher_program_unique'
          AND conrelid = 'public.teacher_programs'::regclass
          AND contype = 'u'
    ) THEN
        ALTER TABLE public.teacher_programs
            ADD CONSTRAINT teacher_program_unique UNIQUE (teacher_id, program_id);
    END IF;
END
$$;

COMMIT;
