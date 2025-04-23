BEGIN;

-- 1. Удалить индексы на таблице referral
DROP INDEX IF EXISTS idx_referral_referrer_id;
DROP INDEX IF EXISTS idx_referral_referee_id;

-- 2. Удалить саму таблицу referral
DROP TABLE IF EXISTS referral;

-- 3. Удалить уникальный индекс по telegram_id
ALTER TABLE customer
DROP CONSTRAINT IF EXISTS customer_telegram_id_unique;

COMMIT;
