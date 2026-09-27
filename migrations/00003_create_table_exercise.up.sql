create table if not exists "exercises" (
    "id" uuid primary key default gen_random_uuid(),
    "exercise_name" varchar(100) unique not null,
    "target_muscle" varchar(100),
    "created_at" timestamptz default now()
);
