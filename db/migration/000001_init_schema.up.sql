CREATE TABLE "public"."users" (
                                  "id" BIGSERIAL PRIMARY KEY,
                                  "phone" VARCHAR NOT NULL UNIQUE,
                                  "email" VARCHAR NOT NULL UNIQUE,
                                  "username" VARCHAR NOT NULL UNIQUE,
                                  "password" VARCHAR NOT NULL,
                                  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
                                  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
                                  "deleted_at" TIMESTAMPTZ
);

-- 重新创建唯一索引，防止重复
CREATE UNIQUE INDEX "idx_phone" ON "public"."users" ("phone");
CREATE UNIQUE INDEX "idx_email" ON "public"."users" ("email");
CREATE UNIQUE INDEX "idx_username" ON "public"."users" ("username");
