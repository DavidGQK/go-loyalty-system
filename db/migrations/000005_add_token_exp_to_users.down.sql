BEGIN TRANSACTION;

ALTER TABLE users
DROP COLUMN IF EXISTS token_exp_at;

COMMIT;
