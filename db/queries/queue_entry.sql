-- name: CreateQueueEntry :one
INSERT INTO queue_entries (room_id, user_id, from_stop_id, to_stop_id, time_start, time_end, min_seats, max_seats)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetQueueEntryByID :one
SELECT * FROM queue_entries WHERE id = $1;

-- name: GetQueueEntryByRoomID :one
SELECT
  qe.id,
  qe.room_id,
  qe.user_id,
  qe.from_stop_id,
  qe.to_stop_id,
  qe.time_start,
  qe.time_end,
  qe.min_seats,
  qe.max_seats,
  qe.status,
  qe.created_at,
  mm.match_id
FROM queue_entries qe
LEFT JOIN match_members mm ON mm.queue_entry_id = qe.id AND mm.user_id = qe.user_id
WHERE qe.id = $1 AND qe.room_id = $2;

-- name: ListWaitingQueueEntriesByRoomID :many
SELECT * FROM queue_entries
WHERE room_id = $1 AND status = 'waiting'
ORDER BY created_at;

-- name: ListMyActiveQueueEntries :many
SELECT
  qe.id,
  qe.room_id,
  qe.user_id,
  qe.from_stop_id,
  qe.to_stop_id,
  qe.time_start,
  qe.time_end,
  qe.min_seats,
  qe.max_seats,
  qe.status,
  qe.created_at,
  mm.match_id
FROM queue_entries qe
LEFT JOIN match_members mm ON mm.queue_entry_id = qe.id AND mm.user_id = qe.user_id
WHERE qe.user_id = $1
  AND qe.status IN ('waiting', 'matched')
ORDER BY qe.created_at DESC;

-- name: CancelQueueEntry :one
UPDATE queue_entries
SET status = 'cancelled'
WHERE id = $1 AND room_id = $2 AND user_id = $3 AND status = 'waiting'
RETURNING *;
