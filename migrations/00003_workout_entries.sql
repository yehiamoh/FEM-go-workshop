-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workout_entries (
   id BIGSERIAL PRIMARY KEY,
   workout_id BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,  -- Fixed typo: REFRENCES -> REFERENCES
   exercise_name VARCHAR(255) NOT NULL,  -- Added space for readability
   sets INTEGER NOT NULL,
   reps INTEGER NOT NULL,
   duration_seconds INTEGER,
   weight DECIMAL(5,2),
   notes TEXT,
   order_index INTEGER NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   CONSTRAINT valid_workout_entry CHECK (
      (reps IS NOT NULL OR duration_seconds IS NOT NULL) AND
      (reps IS NULL OR duration_seconds IS NULL)
   )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workout_entries;
-- +goose StatementEnd