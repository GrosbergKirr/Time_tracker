-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
                                     id serial primary key,
                                     name varchar,
                                     surname varchar,
                                     patronymic varchar,
                                     address varchar,
                                     passport_number varchar
);
CREATE TABLE IF NOT EXISTS tasks(
                                    id serial primary key,
                                    name varchar,
                                    time_begin timestamp default current_timestamp,
                                    time_end timestamp,
                                    user_id int,
                                    status varchar,
                                    foreign key (user_id) references users (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_users_surname ON users(name);
CREATE INDEX IF NOT EXISTS idx_users_surname ON users(surname);
CREATE INDEX IF NOT EXISTS idx_users_passport_number ON users(passport_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop table users, tasks
-- +goose StatementEnd
