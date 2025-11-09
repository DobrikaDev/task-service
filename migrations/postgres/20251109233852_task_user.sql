-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_tasks (
    user_id VARCHAR(255) NOT NULL REFERENCES users(id),
    task_id VARCHAR(255) NOT NULL REFERENCES tasks(id),
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, task_id)
);

CREATE INDEX IF NOT EXISTS idx_user_tasks_user_id ON user_tasks (user_id);
CREATE INDEX IF NOT EXISTS idx_user_tasks_task_id ON user_tasks (task_id);
CREATE INDEX IF NOT EXISTS idx_user_tasks_created_at ON user_tasks (created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_tasks_created_at;
DROP INDEX IF EXISTS idx_user_tasks_task_id;
DROP INDEX IF EXISTS idx_user_tasks_user_id;
DROP TABLE IF EXISTS user_tasks;
-- +goose StatementEnd
