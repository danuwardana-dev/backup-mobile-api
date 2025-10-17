CREATE TABLE IF NOT EXISTS access_states (
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    id bigserial not null primary key,
    user_id bigint not null
        constraint fk_user_id_reset_pin
            references users (id),
    user_uuid varchar(36) not null,
    device_id varchar(50) not null,
    access_token varchar(50) not null,
    access_type varchar(50),
    used boolean not null,
    expired_at timestamp with time zone not null
);
