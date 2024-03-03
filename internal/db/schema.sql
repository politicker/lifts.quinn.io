create table
    lift_set_log
(
    workout_name     text      not null,
    workout_duration text      not null,
    exercise_name    text      not null,
    weight           float     not null,
    reps             float     not null,
    distance         float     not null,
    seconds          float     not null,

    -- Logged at is when the workout actually took place. Imported is when we pulled it into the DB
    logged_at        timestamp,
    imported_at      timestamp not null default current_timestamp
);
