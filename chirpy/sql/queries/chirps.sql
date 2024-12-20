-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid (), NOW(), NOW(), $1, $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * from chirps
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT * from chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE from chirps
WHERE id = $1;

-- name: GetChirpsByUserId :many
SELECT * from chirps
WHERE user_id = $1
ORDER BY created_at ASC;