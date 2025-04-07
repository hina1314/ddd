-- name: GetUserByID :one
SELECT * FROM "user" WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM "user" WHERE username = $1;

-- name: GetUserByPhone :one
SELECT * FROM "user" WHERE phone = $1;

-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email = $1;

-- name: CreateUser :one
INSERT INTO "user" (phone, email, username, password)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: UpdateUser :exec
UPDATE "user" SET phone = $2, email = $3, username = $4, password = $5
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM "user" WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM "user" LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM "user";