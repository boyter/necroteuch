-- name: GetDataByKey :one
SELECT uid, key, value FROM data
WHERE key = ? LIMIT 1;
