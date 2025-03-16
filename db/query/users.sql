-- name: CreateUser :one
INSERT INTO users (
    phone,
    username,
    password
) VALUES (
             $1, $2, $3
         )
    RETURNING *;

-- name: FindUserByPhone :one
SELECT * FROM users WHERE phone = $1;

-- name: UpdateUserInfo :execresult
UPDATE users SET username = $2 where id = $1;

-- name: UpdateUserPassword :execresult
UPDATE users SET password = $2 where id = $1;

-- name: DeleteAuthor :exec
DELETE FROM users
WHERE id = $1;