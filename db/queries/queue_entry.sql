-- name: CreateQueueEntry :one
INSERT INTO queue_entries (room_id, user_id, from_stop_id, to_stop_id, time_start, time_end, min_seats, max_seats)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetQueueEntryByID :one
SELECT * FROM queue_entries WHERE id = $1;

-- name: ListWaitingQueueEntriesByRoomID :many
SELECT * FROM queue_entries
WHERE room_id = $1 AND status = 'waiting'
ORDER BY created_at;

-- name: CancelQueueEntry :one
UPDATE queue_entries
SET status = 'cancelled'
WHERE id = $1 AND room_id = $2 AND user_id = $3 AND status = 'waiting'
RETURNING *;
