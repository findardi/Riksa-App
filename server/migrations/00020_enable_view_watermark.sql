-- +goose Up
-- +goose StatementBegin
update access_levels
set can_watermark = true,
    updated_at = now()
where workspace_id is null and name = 'view';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
update access_levels
set can_watermark = false,
    updated_at = now()
where workspace_id is null and name = 'view';
-- +goose StatementEnd
