-- +goose Up
CREATE TABLE queue_entries (
  id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id      UUID        NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  from_stop_id UUID        NOT NULL REFERENCES hub_stops(id),
  to_stop_id   UUID        NOT NULL REFERENCES hub_stops(id),
  time_start   TIMESTAMPTZ NOT NULL,
  time_end     TIMESTAMPTZ NOT NULL,
  min_seats    INT         NOT NULL DEFAULT 2,
  max_seats    INT         NOT NULL DEFAULT 4,
  status       TEXT        NOT NULL DEFAULT 'waiting',
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_queue_entries_room_id ON queue_entries(room_id);
CREATE INDEX idx_queue_entries_status ON queue_entries(room_id, status);

-- +goose Down
DROP TABLE IF EXISTS queue_entries;
