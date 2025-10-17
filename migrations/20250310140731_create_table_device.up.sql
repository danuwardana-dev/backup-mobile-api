CREATE TABLE IF NOT EXISTS devices(
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    id bigserial not null primary key,
    device_id varchar(50) not null,
    user_id bigint not null
    constraint fk_user_id
    references users (id),
    user_uuid varchar(36) not null,
    app_version_code varchar(50) not null,
    app_version_name varchar(50) not null,
    manufacturer varchar(50) not null,
    brand varchar(50) not null,
    device_model varchar(50) not null,
    product varchar(50) not null,
    version_sdk varchar(50) not null,
    version_release varchar(50) not null

);


