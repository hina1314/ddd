-- name: CreateAccount :one
INSERT INTO user_account (user_id, frozen_balance, balance)
VALUES ($1, $2, $3)
    RETURNING *;