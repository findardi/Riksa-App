-- +goose Up
-- +goose StatementBegin
create table if not exists workspace_user_invitations(
  id uuid primary key default gen_random_uuid(),
  workspace_id uuid not null references workspaces (id) on delete cascade,
  email text not null,
  role_id uuid not null references workspace_roles (id) on delete restrict,
  user_id uuid references users (id) on delete cascade,
  invited_by uuid not null references users (id) on delete cascade,
  code_hash text not null,
  status text not null default 'pending',
  expires_at timestamptz not null,
  accepted_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  constraint workspace_invitations_status_check check (status in ('pending', 'accepted', 'rejected', 'revoked', 'expired'))
);

create unique index workspace_invitations_code_hash_key on workspace_user_invitations (code_hash);

create unique index workspace_invitations_pending_key on workspace_user_invitations (workspace_id, lower(email)) where status = 'pending';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists workspace_user_invitations;
-- +goose StatementEnd
