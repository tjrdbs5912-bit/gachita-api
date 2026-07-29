-- name: AddRoomMember :one
INSERT INTO room_members (room_id, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetRoomMember :one
SELECT * FROM room_members
WHERE room_id = $1 AND user_id = $2;

-- name: ListRoomMembersByRoomID :many
SELECT * FROM room_members WHERE room_id = $1;

-- name: ListRoomMembersWithUser :many
SELECT
  rm.room_id,
  rm.user_id,
  rm.joined_at,
  u.nickname,
  u.email
FROM room_members rm
JOIN users u ON u.id = rm.user_id
WHERE rm.room_id = $1
ORDER BY rm.joined_at;