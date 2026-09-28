-- name: GetUserByID :one
SELECT id, phone, email, username, password, avatar, created_at, updated_at, deleted_at
FROM "user"
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT id, phone, email, username, password, avatar, created_at, updated_at, deleted_at
FROM "user"
WHERE username = $1 AND deleted_at IS NULL;

-- name: GetUserByPhone :one
SELECT id, phone, email, username, password, avatar, created_at, updated_at, deleted_at
FROM "user"
WHERE phone = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT id, phone, email, username, password, avatar, created_at, updated_at, deleted_at
FROM "user"
WHERE email = $1 AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO "user" (phone, email, username, password)
VALUES ($1, $2, $3, $4)
RETURNING id, phone, email, username, password, avatar, created_at, updated_at, deleted_at;

-- name: UpdateUser :execrows
UPDATE "user"
SET phone = $2, email = $3, username = $4, password = $5, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUser :execrows
UPDATE "user" SET deleted_at = now(), updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT id, phone, email, username, password, avatar, created_at, updated_at, deleted_at
FROM "user"
WHERE deleted_at IS NULL AND id > sqlc.arg(after_id)
ORDER BY id
LIMIT sqlc.arg(page_size);
