CREATE TABLE IF NOT EXISTS users (
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    id bigserial not null primary key,
    uuid varchar(36) not null,
    full_name varchar(255) not null,
    roles varchar(20),
    email varchar(255) not null,
    phone_number varchar(15) not null,
    password varchar(255) not null,
    pin varchar(60),
    device_id varchar(20) not null,
    status varchar(50) not null
);
