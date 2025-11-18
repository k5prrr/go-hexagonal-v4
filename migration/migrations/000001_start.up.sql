/*
Позже ещё подумаю, как оптимизировать
user_group_managers
user_group_sellers

*/


BEGIN;



-- Пользователи
CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."users" (
    "id" bigint DEFAULT nextval('users_id_seq') NOT NULL,

    "family_name" character varying(128),
    "name" character varying(128),
    "middle_name" character varying(128),

    "phone" character varying(32) UNIQUE NOT NULL,
    "email" character varying(64),

    "birth_date" date,
    "gender_id" integer,

    "parent_id" bigint,
    "role_id" bigint,

    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW(),
    "deleted_at" timestamptz,
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_phone ON users (phone);
CREATE INDEX idx_users_role_id ON users (role_id);
COMMENT ON TABLE "public"."users" IS 'Пользователи системы';


CREATE SEQUENCE auth_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE auth (
    id bigint DEFAULT nextval('auth_id_seq') NOT NULL,
    user_id bigint NOT NULL,

    tg_id bigint UNIQUE NOT NULL,
    code character varying(16),
    token character varying(64) NOT NULL,

    last_login_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    deleted_at timestamptz
);
COMMENT ON TABLE "public"."auth" IS 'Авторизация и Аутентификация';




CREATE SEQUENCE user_roles_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."user_roles" (
                                        "id" bigint DEFAULT nextval('user_roles_id_seq') NOT NULL,
                                        "name" character varying(16) NOT NULL,
                                        CONSTRAINT "user_roles_pkey" PRIMARY KEY ("id")
);
INSERT INTO "user_roles" ("id", "name") VALUES
                                             (1,	'Админ'),
                                             (2,	'Менеджер'),
                                             (3,	'Продавец'),
                                             (4,	'Новый');
COMMENT ON TABLE "public"."user_roles" IS 'Должности/Роли пользователей системы';




CREATE SEQUENCE user_group_managers_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."user_group_managers" (
                                                "id" bigint DEFAULT nextval('user_group_managers_id_seq') NOT NULL,
                                                "parent_id" bigint,
                                                "user_id" bigint,
                                                "shop_group_id" integer,
                                                CONSTRAINT "user_group_managers_pkey" PRIMARY KEY ("id")
);
COMMENT ON TABLE "public"."user_group_managers" IS 'Поля менеджеров';




CREATE SEQUENCE user_group_sellers_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."user_group_sellers" (
    "id" bigint DEFAULT nextval('user_group_sellers_id_seq') NOT NULL,
    "parent_id" bigint,
    "user_id" bigint,
    "shop_group_id" integer,
    "shop_id" bigint,
    CONSTRAINT "user_group_sellers_pkey" PRIMARY KEY ("id")
);
COMMENT ON TABLE "public"."user_group_sellers" IS 'Поля продавцов';


CREATE SEQUENCE auth_codes_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE auth_codes (
    "id" bigint DEFAULT nextval('auth_codes_id_seq') NOT NULL,
    "code" character varying(128) NOT NULL,
    "type" character varying(16) NOT NULL,
    "uuid" character varying(64),
    "phone" character varying(64),
    "created_at" timestamptz NOT NULL,
    "updated_at" timestamptz NOT NULL,
    CONSTRAINT "auth_codes_pkey" PRIMARY KEY ("id")
);




-- Клиенты


CREATE SEQUENCE clients_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."clients" (
                                    "id" bigint DEFAULT nextval('clients_id_seq') NOT NULL,
                                    "family_name" character varying(128),
                                    "name" character varying(128),
                                    "middle_name" character varying(128),
                                    "phone" character varying(32),
                                    "email" character varying(64),
                                    "birth_date" date,
                                    "age" integer,
                                    "parent_id" bigint,
                                    "gender_id" integer,
                                    "shop_id" character varying(128),
                                    "product_code" character varying(64),
                                    "product_name" character varying(128),
                                    "old_product_id" integer,
                                    "old_product_name" character varying(64),
                                    "old_product_model" character varying(64),
                                    "quantity_cartridge06" integer,
                                    "quantity_cartridge08" integer,
                                    "quantity_cartridge1" integer,
                                    "created_at" timestamptz DEFAULT NOW(),
                                    "updated_at" timestamptz DEFAULT NOW(),
                                    "deleted_at" timestamp,
                                    CONSTRAINT "clients_pkey" PRIMARY KEY ("id")
);
CREATE INDEX idx_clients_email ON clients (email);
CREATE INDEX idx_clients_phone ON clients (phone);
COMMENT ON TABLE "public"."clients" IS 'Клиенты';





-- Возвраты



CREATE SEQUENCE refunds_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."refunds" (
                                    "id" bigint DEFAULT nextval('refunds_id_seq') NOT NULL,
                                    "clients_id" bigint,
                                    "parent_id" bigint NOT NULL,
                                    "shop_id" character varying(128),
                                    "reason" text,
                                    "product_code" character varying(64),
                                    "product_name" character varying(128),
                                    "created_at" timestamptz DEFAULT NOW(),
                                    "updated_at" timestamptz DEFAULT NOW(),
                                    "deleted_at" timestamp,
                                    CONSTRAINT "refunds_pkey" PRIMARY KEY ("id")
);
COMMENT ON COLUMN "public"."refunds"."reason" IS 'Причина возврата';
COMMENT ON TABLE "public"."refunds" IS 'Возвраты';






-- Магазины

CREATE SEQUENCE shop_roles_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;
CREATE TABLE "public"."shop_roles" (
                                        "id" integer DEFAULT nextval('shop_roles_id_seq') NOT NULL,
                                        "name" character varying(32) NOT NULL,
                                        CONSTRAINT "shop_roles_pkey" PRIMARY KEY ("id")
);
COMMENT ON TABLE "public"."shop_roles" IS 'Сети магазинов';




CREATE SEQUENCE shops_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."shops" (
                                  "id" bigint DEFAULT nextval('shops_id_seq') NOT NULL,
                                  "parent_id" bigint NOT NULL,
                                  "address" character varying(256) NOT NULL,
                                  "group_id" integer NOT NULL,
                                  "city_id" integer NOT NULL,
                                  CONSTRAINT "shops_pkey" PRIMARY KEY ("id")
);
COMMENT ON TABLE "public"."shops" IS 'Магазины';






-- Поставки


CREATE SEQUENCE postings_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."postings" (
                                     "id" bigint DEFAULT nextval('postings_id_seq') NOT NULL,
                                     "parent_id" bigint NOT NULL,
                                     "shop_id" bigint NOT NULL,
                                     "description" text NOT NULL,
                                     "quantity" integer NOT NULL,
                                     "status_id" integer NOT NULL,
                                     "created_at" timestamptz DEFAULT NOW(),
                                     "updated_at" timestamptz DEFAULT NOW(),
                                     "deleted_at" timestamp,
                                     CONSTRAINT "postings_pkey" PRIMARY KEY ("id")
);
CREATE INDEX idx_postings_shop_id ON postings (shop_id);
CREATE INDEX idx_postings_status_id ON postings (status_id);
CREATE INDEX idx_postings_created_at ON postings (created_at);
COMMENT ON COLUMN "public"."postings"."quantity" IS 'Количество единиц товара';
COMMENT ON TABLE "public"."postings" IS 'Поставки товаров в магазины';





CREATE SEQUENCE posting_status_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;
CREATE TABLE "public"."posting_status" (
                                           "id" integer DEFAULT nextval('posting_status_id_seq') NOT NULL,
                                           "name" character varying(64) NOT NULL,
                                           "sort" integer,
                                           CONSTRAINT "posting_status_pkey" PRIMARY KEY ("id")
);
INSERT INTO "posting_status" ("id", "name", "sort") VALUES
                                                        (1,	'Создана', 10),
                                                        (2,	'Ожидает отправки', 20),
                                                        (3,	'Отправлена', 30),
                                                        (4,	'В пути', 40),
                                                        (5,	'Доставлена', 50),
                                                        (6,	'Отклонена', 60),
                                                        (7,	'Отменена', 70),
                                                        (8,	'Проблема с доставкой', 80),
                                                        (9,	'Возвращена', 90);
CREATE INDEX idx_posting_status_name ON posting_status (name);
CREATE INDEX idx_posting_status_sort ON posting_status (sort);
COMMENT ON TABLE "public"."posting_status" IS 'Статусы поставок';







-- Прочее


CREATE SEQUENCE old_products_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;
CREATE TABLE "public"."old_products" (
                                         "id" integer DEFAULT nextval('old_products_id_seq') NOT NULL,
                                         "name" character varying(32) NOT NULL,
                                         "created_at" timestamptz DEFAULT NOW(),
                                         "updated_at" timestamptz DEFAULT NOW(),
                                         "deleted_at" timestamp,
                                         CONSTRAINT "old_products_pkey" PRIMARY KEY ("id")
);
INSERT INTO "old_products" ("id", "name") VALUES
                                              (1,	'Другой'),
                                              (63, 'VAPORESSO'),
                                              (64, 'BRUSKO'),
                                              (65, 'VOOPOO'),
                                              (66, 'GEEKVAPE'),
                                              (67, 'OXVA'),
                                              (68, 'SMOANT'),
                                              (69, 'RINCOE'),
                                              (70, 'LOST VAPE');
COMMENT ON TABLE "public"."old_products" IS 'Бренды старых устройств';




CREATE SEQUENCE genders_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;
CREATE TABLE "public"."genders" (
                                    "id" integer DEFAULT nextval('genders_id_seq') NOT NULL,
                                    "name" character varying(16) NOT NULL,
                                    CONSTRAINT "genders_pkey" PRIMARY KEY ("id")
);
INSERT INTO "genders" ("id", "name") VALUES
                                         (1,	'Мужской'),
                                         (2,	'Женский');
COMMENT ON TABLE "public"."genders" IS 'Пол';





COMMIT;





