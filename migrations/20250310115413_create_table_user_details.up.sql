CREATE TABLE IF NOT EXISTS user_details (
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    id bigserial not null primary key,
    user_id bigint not null
    constraint fk_user_id
    references users (id),
    user_uuid varchar(36) not null,
    country varchar(50),
    province varchar(50),
    regency varchar(50),
    district varchar(50),
    address varchar (300),
    biometric varchar(50),
    status varchar(50),
    kyc_status varchar(50) not null,
    kyc_type varchar(20),
    profile_picture varchar(250)
);
