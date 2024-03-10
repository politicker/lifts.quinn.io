-- name: Get1RMHistory :many
WITH OrderedSets AS (SELECT DISTINCT ON (logged_at) logged_at,
                                                    weight,
                                                    reps,
                                                    (weight * reps * 0.0333 + weight)::float AS estimated_1rm
                     FROM lift_set_log
                     WHERE lower(exercise_name) = lower($1)
                     ORDER BY logged_at DESC, estimated_1rm DESC
                     LIMIT 11)
SELECT estimated_1rm, logged_at, (round(weight)::text || ' x ' || reps::text) as set_text
FROM OrderedSets
ORDER BY logged_at ASC;

-- name: GetBestSet :one
SELECT (round(weight)::text || ' x ' || reps::text) AS set_text,
       weight,
       reps,
       logged_at,
       (weight * reps * 0.0333 + weight)            as estimated_1rm
FROM lift_set_log
WHERE LOWER(exercise_name) = LOWER($1)
ORDER BY estimated_1rm DESC,
         logged_at DESC
LIMIT 1;

-- name: CreateLiftSetLog :exec
INSERT INTO lift_set_log (workout_name, workout_duration, exercise_name, set_order, weight, reps, distance, seconds,
                          notes, workout_notes, rpe, logged_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT (workout_name, exercise_name, set_order, logged_at) DO NOTHING;
