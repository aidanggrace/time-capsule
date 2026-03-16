-- name: ListCapsules :many
SELECT * FROM capsules
WHERE owner_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;


-- name: CreateCapsule :one
INSERT INTO capsules (
  owner_id,
  recipient_email,
  message,
  media_url,
  unlock_at
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetCapsule :one
SELECT * FROM capsules
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  email,
  password_hash
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

