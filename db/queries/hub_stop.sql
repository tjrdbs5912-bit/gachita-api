-- name: CreateHubStop :one
INSERT INTO hub_stops (room_id, name, sort_order)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListHubStopsByRoomID :many
SELECT * FROM hub_stops
WHERE room_id = $1
ORDER BY sort_order, created_at;

-- name: GetHubStopByID :one
SELECT * FROM hub_stops WHERE id = $1;

-- name: UpdateHubStop :one
UPDATE hub_stops
SET name = $3, sort_order = $4
WHERE id = $1 AND room_id = $2
RETURNING *;

-- name: DeleteHubStop :one
DELETE FROM hub_stops
WHERE id = $1 AND room_id = $2
RETURNING *;