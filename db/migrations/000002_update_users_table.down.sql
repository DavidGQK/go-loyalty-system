BEGIN TRANSACTION;

ALTER TABLE users
DROP COLUMN bonuses;

COMMIT;
