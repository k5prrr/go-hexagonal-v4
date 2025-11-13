
BEGIN;


DROP TABLE IF EXISTS "users";
DROP SEQUENCE IF EXISTS users_id_seq;

DROP TABLE IF EXISTS "user_groups";
DROP SEQUENCE IF EXISTS user_groups_id_seq;

DROP TABLE IF EXISTS "user_group_managers";
DROP SEQUENCE IF EXISTS user_group_managers_id_seq;

DROP TABLE IF EXISTS "user_group_sellers";
DROP SEQUENCE IF EXISTS user_group_sellers_id_seq;

DROP TABLE IF EXISTS "auth_codes";
DROP SEQUENCE IF EXISTS auth_codes_id_seq;




-- Клиенты
DROP TABLE IF EXISTS "clients";
DROP SEQUENCE IF EXISTS clients_id_seq;

-- Возвраты
DROP TABLE IF EXISTS "refunds";
DROP SEQUENCE IF EXISTS refunds_id_seq;

-- Магазины
DROP TABLE IF EXISTS "shop_groups";
DROP SEQUENCE IF EXISTS shop_groups_id_seq;

DROP TABLE IF EXISTS "shops";
DROP SEQUENCE IF EXISTS shops_id_seq;

-- Поставки
DROP TABLE IF EXISTS "postings";
DROP SEQUENCE IF EXISTS postings_id_seq;

DROP TABLE IF EXISTS "posting_status";
DROP SEQUENCE IF EXISTS posting_status_id_seq;

-- Прочее
DROP TABLE IF EXISTS "old_products";
DROP SEQUENCE IF EXISTS old_products_id_seq;


DROP TABLE IF EXISTS "genders";
DROP SEQUENCE IF EXISTS genders_id_seq;



COMMIT;

