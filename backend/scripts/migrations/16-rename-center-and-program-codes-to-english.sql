-- 16-rename-center-and-program-codes-to-english.sql
-- Normalize educational center and program codes to the new standard:
-- educational_centers.code: cepprs
-- programs.code: seeds, fly, cecor
--
-- This migration is idempotent.

BEGIN;

-- 1) Educational center code
UPDATE public.educational_centers
SET code = 'cepprs',
    updated_at = NOW()
WHERE code IN ('CEPROS', 'cepros')
   OR (name = 'Centro Educacional Prof. Paulo Rossi Severino' AND code <> 'cepprs');

-- 2) Program codes
UPDATE public.programs
SET code = 'seeds',
    updated_at = NOW()
WHERE code IN ('SEMEAR', 'semear')
   OR (name = 'Semear' AND code <> 'seeds');

UPDATE public.programs
SET code = 'fly',
    updated_at = NOW()
WHERE code IN ('VOAR', 'voar')
   OR (name = 'Voar' AND code <> 'fly');

UPDATE public.programs
SET code = 'cecor',
    updated_at = NOW()
WHERE code IN ('CECOR')
   OR (name = 'Cecor' AND code <> 'cecor');

COMMIT;
