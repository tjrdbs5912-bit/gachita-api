-- name: CreateUser :one
INSERT INTO users (email, password_hash, nickname)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT 
    id,
    email,
    nickname,
    password_hash,
    created_at,
    updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT 
    id,
    email,
    nickname,
    password_hash,
    created_at,
    updated_at
FROM users
WHERE email = $1;

-- name: UpdateUserProfile :one
UPDATE users
SET
  email = $2,
  nickname = $3,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET
  password_hash = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;