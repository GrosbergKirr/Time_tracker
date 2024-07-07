-- +goose Up
-- +goose StatementBegin
alter table tasks add status varchar;
-- +goose StatementEnd


