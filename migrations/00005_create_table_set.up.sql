create table if not exists "sets" (
    "id" uuid primary key default gen_random_uuid(),
    "set_number" int not null,
    "weight" numeric(10,2) default 0,
    "reps" int default 0,
    "time" int default 0,
    "created_at" timestamptz default now(),
    "workout_exercises_id" uuid not null references workout_exercises(id) on delete cascade
);

create index if not exists idx_sets_workout_exercises_id on sets(workout_exercises_id);
