create table if not exists "workout_exercises" (
    "id" uuid primary key default gen_random_uuid(),
    "exercise_id" uuid not null references exercises(id) on delete restrict,
    "workout_id" uuid not null references workouts(id) on delete cascade,
    "created_at" timestamptz default now()
);

create index if not exists idx_workout_exercises_workout_id on workout_exercises(workout_id);
create index if not exists idx_workout_exercises_exercise_id on workout_exercises(exercise_id);
