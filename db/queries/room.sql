-- name: CreateRoom :one
INSERT INTO rooms (name, invite_code, openchat_url, owner_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRoomByID :one
SELECT * FROM rooms WHERE id = $1;

-- name: GetRoomByInviteCode :one
SELECT * FROM rooms WHERE invite_code = $1;

-- name: ListRoomsByUserID :many
SELECT r.*
FROM rooms r
JOIN room_members rm ON rm.room_id = r.id
WHERE rm.user_id = $1
ORDER BY rm.joined_at DESC;