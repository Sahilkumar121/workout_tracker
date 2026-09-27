-- 1. Insert 30 dummy users and return their generated UUIDs
WITH inserted_users AS (
    INSERT INTO users (first_name, last_name, email, password, created_at)
    SELECT
        (ARRAY[
            'Alex', 'Jordan', 'Taylor', 'Morgan', 'Sam', 'Chris', 'Pat', 'Casey', 'Riley', 'Dakota',
            'Avery', 'Reese', 'Rowan', 'Quinn', 'Skyler', 'Cameron', 'Jesse', 'Peyton', 'Logan', 'Harper',
            'Elliott', 'Emerson', 'Finley', 'Hayden', 'Kai', 'Parker', 'Sage', 'Shiloh', 'Tatum', 'Val'
        ])[i] AS first_name,
        (ARRAY[
            'Smith', 'Johnson', 'Williams', 'Brown', 'Jones', 'Garcia', 'Miller', 'Davis', 'Rodriguez', 'Martinez',
            'Hernandez', 'Lopez', 'Gonzalez', 'Wilson', 'Anderson', 'Thomas', 'Taylor', 'Moore', 'Jackson', 'Martin',
            'Lee', 'Perez', 'Thompson', 'White', 'Harris', 'Sanchez', 'Clark', 'Ramirez', 'Lewis', 'Robinson'
        ])[i] AS last_name,
        LOWER((ARRAY[
            'Alex', 'Jordan', 'Taylor', 'Morgan', 'Sam', 'Chris', 'Pat', 'Casey', 'Riley', 'Dakota',
            'Avery', 'Reese', 'Rowan', 'Quinn', 'Skyler', 'Cameron', 'Jesse', 'Peyton', 'Logan', 'Harper',
            'Elliott', 'Emerson', 'Finley', 'Hayden', 'Kai', 'Parker', 'Sage', 'Shiloh', 'Tatum', 'Val'
        ])[i]) || i || '@example.com' AS email,
        
        -- Valid bcrypt hash for the password 'password123456' (14 characters)
        -- Satisfies min=12 validation constraint and passes bcrypt comparison.
        '$2a$12$3f0aPeMwAcvicvlxhdiuW.2hIiUQUK28EUbajPh5Fy/hKQ8uJI6bu' AS password,
        
        NOW() - (i || ' days')::INTERVAL AS created_at
    FROM generate_series(1, 30) AS i
    RETURNING id
)

-- 2. Insert 30 dummy workouts linked to the created users
INSERT INTO workouts (workout_name, workout_start, created_at, user_id)
SELECT
    (ARRAY[
        'Morning Cardio', 'Leg Day Blast', 'Upper Body Power', 'Full Body HIIT',
        'Yoga Flow', 'Core Crusher', 'Evening Run', 'Push Day',
        'Pull Day', 'Swimming Session', 'CrossFit Metcon', 'Chest & Triceps',
        'Back & Biceps', 'Pilates Rebound', 'Cycling Sprint', 'Functional Mobility'
    ])[floor(random() * 16 + 1)] AS workout_name,
    NOW() - (floor(random() * 20) || ' days')::INTERVAL - (floor(random() * 12) || ' hours')::INTERVAL AS workout_start,
    NOW() - (floor(random() * 30) || ' days')::INTERVAL AS created_at,
    id AS user_id
FROM inserted_users;