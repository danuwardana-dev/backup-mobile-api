CREATE TABLE IF NOT EXISTS otps(
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    id bigserial not null primary key,
    otp varchar(6) not null,
    otp_purpose varchar(255) not null,
    otp_method varchar(50) not null,
    otp_destination varchar(50) not null,
    user_id bigint not null
    constraint fk_user_id_reset_pin
    references users (id),
    user_uuid varchar(36) not null,
    identity_user varchar(255) not null,
    verify_key varchar(50) not null,
    session_id varchar(50) not null,
    status  varchar(50) not null,
    expired_at timestamp with time zone not null
);

