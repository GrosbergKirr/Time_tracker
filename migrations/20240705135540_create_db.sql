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
                                    time_end timestamp default current_timestamp,
                                    user_id int,
                                    foreign key (user_id) references users (id) ON DELETE CASCADE
);
-- +goose StatementEnd




