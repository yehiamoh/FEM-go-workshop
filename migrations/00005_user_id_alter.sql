-- +goose Up
-- +goose StatementBegin
ALTER table workouts
ADD Column user_id BIGINT not null REFERENCES users(id) on DELETE CASCADE
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workouts DROP Column user_id;
-- +goose StatementEnd