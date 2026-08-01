-- name: ListMatchCandidates :many
SELECT * FROM queue_entries
WHERE room_id = $1
  AND from_stop_id = $2
  AND to_stop_id = $3
  AND status = 'waiting'
  AND time_start < $4
  AND time_end > $5
ORDER BY created_at;

-- name: CreateMatch :one
INSERT INTO matches (room_id)
VALUES ($1)
RETURNING *;

-- name: AddMatchMember :exec
INSERT INTO match_members (match_id, user_id, queue_entry_id)
VALUES ($1, $2, $3);

-- name: MarkQueueEntriesMatched :exec
UPDATE queue_entries
SET status = 'matched'
WHERE id = ANY($1::uuid[]);

-- name: ListMatchesByRoomID :many
SELECT * FROM matches
WHERE room_id = $1
ORDER BY created_at DESC;

-- name: GetMatchByID :one
SELECT * FROM matches WHERE id = $1;

-- name: ListMatchMembersByMatchID :many
SELECT
  mm.match_id,
  mm.user_id,
  mm.queue_entry_id,
  u.nickname,
  qe.from_stop_id,
  qe.to_stop_id,
  qe.time_start,
  qe.time_end
FROM match_members mm
JOIN users u ON u.id = mm.user_id
JOIN queue_entries qe ON qe.id = mm.queue_entry_id
WHERE mm.match_id = $1;
