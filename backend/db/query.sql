-- name: ListCapsules :many
SELECT * FROM capsules
WHERE deleted_at IS NULL
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
