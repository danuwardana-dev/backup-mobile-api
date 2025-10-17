CREATE TABLE IF NOT EXISTS identity_ktps (
    id bigserial not null primary key,
    user_id bigint not null
        constraint fk_user_id_ktp
            references users (id),
    nik varchar(255) not null,
    full_name varchar(255) not null,
    place_of_birth varchar(100) not null,
    gender varchar(20) not null,
    date_of_birth date not null,
    occupation varchar(100) not null,
    nationality varchar(100) not null,
    martial_status varchar(20) not null,
    religion varchar(20) not null,
    country varchar(100) not null,
    province varchar(100) not null,
    city varchar(100) not null, 
    district varchar(255) not null,
    full_address text not null,
    identity_image varchar(255) not null,
    submit_at timestamp with time zone not null,
    verify_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
);

CREATE TABLE IF NOT EXISTS identity_passports (
    id bigserial not null primary key,
    user_id bigint not null
        constraint fk_user_id_passport
            references users (id),
    passport_number varchar(255) not null,
    passport_type varchar(50) not null,
    gender varchar(20) not null,
    full_name varchar(255) not null,
    nationality varchar(100) not null,
    place_of_birth varchar(100) not null,
    date_of_birth date not null,
    date_of_issue date not null,
    date_of_expired date not null,
    resi_number varchar(255) not null,
    place_of_issue varchar(100) not null,
    identity_image varchar(255) not null,
    submit_at timestamp with time zone not null,
    verify_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
);