CREATE TABLE IF NOT EXISTS articles (
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    id bigserial not null primary key,
    name varchar(250) not null,
    url varchar not null,
    category varchar(50),
    active_after_day timestamp with time zone
);
