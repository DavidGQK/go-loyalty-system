BEGIN TRANSACTION;

ALTER TABLE users
    ADD COLUMN bonuses INT DEFAULT 0;

COMMIT;
