-- name: GetDataByKey :one
SELECT uid, key, value FROM data
WHERE key = ? LIMIT 1;

-- name: DoNothing :exec
SELECT count(*) FROM data;