create table if not exists "users" (
    "id" uuid primary key default gen_random_uuid (),
    "first_name" varchar(50) not null,
    "last_name" varchar(50),
    "email" varchar(100) unique not null,
    "password" varchar(100) not null,
    "created_at" timestamptz default now()
);