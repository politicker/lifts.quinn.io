-- name: Get1RMHistory :many
select distinct on (logged_at) (weight * reps * 0.0333 + weight)::float AS estimated_1rm
from lift_set_log
where lower(exercise_name) = lower($1)
order by logged_at desc,
         estimated_1rm desc
LIMIT 11;

-- name: GetBestSet :one
SELECT (round(weight)::text || ' x ' || reps::text) AS best_set
FROM lift_set_log
WHERE LOWER(exercise_name) = LOWER($1)
ORDER BY logged_at DESC,
         (weight * reps * 0.0333 + weight) DESC
LIMIT 1;

-- name: CreateLiftSetLog :exec
INSERT INTO lift_set_log (name, workout_name, exercise_name, workout_duration, weight, reps, seconds, distance,
                          logged_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
