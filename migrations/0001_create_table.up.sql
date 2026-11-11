
-- USERS : store the data for user
create table if not exists users(
    id uuid primary key default gen_random_uuid(),
    name varchar(100) not null,
    email varchar(100) not null unique,
    password text not null,
    
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- EXERCISES : store all predefined exercise
create table if not exists exercises(
    id uuid primary key default gen_random_uuid(),
    name varchar(20) not null unique,
    description text not null,
    muscle_group varchar(50) not null,
    
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- WORKOUT : store user workout
-- a user can have multiple workout
-- One to many relationship
create table if not exists workouts(
    id uuid primary key default gen_random_uuid(),
    title varchar(50) not null,

    started_at timestamptz not null default now(),
    completed_at timestamptz,

    created_at timestamptz default now(),
    updated_at timestamptz default now(),

    user_id uuid not null,

    constraint fk_workout_user
        foreign key user_id references users(id)
    on delete cascade,

    constraint ck_workout_completed
    check (completed_at is null or completed_at >= started_at)
);

-- WORKOUT EXERCISES : store the workout and exercise
-- many workout can have many exercise
-- many to many relation
create table if not exists workout_exercises(
    id uuid primary key default gen_random_uuid(),

    workout_id uuid not null,
    exercise_id uuid not null,
    positions int not null check (positions > 0),

    created_at timestamptz not null default now()

    constraint fk_workout
        foreign key workout_id references workouts(id)
    on delete cascade,

    constraint fk_exercise
        foreign key exercise_id references exercises(id)
    on delete restrict,

    constraint uq_workout_exercise
        unique(workout_id, exercise_id),


    constraint uq_workout_position
        unique(workout_id, positions)

);

-- WORKOUT SET : store workout_exercises sets, reps, weight
-- each workout_exercises will have many sets
create table if not exists workout_sets(
    id uuid primary key default gen_random_uuid(),
    
    workout_exercises_id uuid not null,
    
    set_number int not null check (set_number > 0),
    reps int not null check (reps > 0),
    weight_kg numeric(7,2) not null check (weight_kg > 0),
    
    created_at timestamptz not null default now()

    constraint fk_workout_exercise
        foreign key workout_exercises_id references workout_exercises(id)
    on delete cascade,

    constraint uq_workout_exercise_set_number
        unique(workout_exercises_id, set_number)
);
