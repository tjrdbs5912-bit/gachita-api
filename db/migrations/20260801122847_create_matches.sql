-- +goose Up
CREATE TABLE matches (
  id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id    UUID        NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  status     TEXT        NOT NULL DEFAULT 'confirmed',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_matches_room_id ON matches(room_id);

CREATE TABLE match_members (
  match_id       UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
  user_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  queue_entry_id UUID NOT NULL REFERENCES queue_entries(id),
  PRIMARY KEY (match_id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS match_members;
DROP TABLE IF EXISTS matches;
