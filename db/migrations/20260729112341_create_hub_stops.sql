-- +goose Up
CREATE TABLE hub_stops (
  id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id    UUID        NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  name       TEXT        NOT NULL,
  sort_order INT         NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_hub_stops_room_id ON hub_stops(room_id);

-- +goose Down
DROP TABLE IF EXISTS hub_stops;