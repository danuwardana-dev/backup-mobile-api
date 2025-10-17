CREATE TABLE IF NOT EXISTS token_blacklists(
    id bigserial not null primary key,
    token varchar  not null,
    blacklist_at timestamp with time zone not null,
    expired_at timestamp with time zone not null,
    description varchar(200),
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);
