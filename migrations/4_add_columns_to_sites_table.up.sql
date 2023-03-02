ALTER TABLE sites
ADD COLUMN IF NOT EXISTS "registrar" json;

ALTER TABLE sites
ADD COLUMN IF NOT EXISTS "registrant" json;

ALTER TABLE sites
ADD COLUMN IF NOT EXISTS "administrative" json;

ALTER TABLE sites
ADD COLUMN IF NOT EXISTS "technical" json;

ALTER TABLE sites
ADD COLUMN IF NOT EXISTS "billing" json;
