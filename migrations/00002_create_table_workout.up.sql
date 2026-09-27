create table if not exists "workouts" (
    "id" uuid primary key default gen_random_uuid(),
    "workout_name" varchar(100) not null,
    "workout_start" timestamptz default now(),
    "created_at" timestamptz default now(),
    "user_id" uuid not null references users(id) on delete cascade
);

create index if not exists idx_workouts_user_id on workouts(user_id);
