-- table: rooms
-- name: InsertRoom :one
INSERT OR IGNORE INTO rooms (name, password) VALUES (?, ?) RETURNING *;
-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id = ?;
-- name: SelectRoomByNameWithPassword :one
SELECT *, password FROM rooms WHERE name = ?;
-- name: SelectRoomById :one
SELECT * FROM rooms WHERE id = ?;
-- name: SelectRoomByName :one
SELECT * FROM rooms WHERE name = ?;
-- name: UpdateRoomNameById :exec
UPDATE rooms SET name = ? WHERE id = ?;
-- name: SelectCurrentVideoByRoomId :one
SELECT video.uri, video.local FROM rooms
JOIN video ON rooms.currently_playing = video.id
WHERE rooms.id = ?;

-- -- table: users
-- -- name: InsertUser :one
-- INSERT INTO users (username) VALUES (?) RETURNING id, username;
-- -- name: DeleteUser :exec
-- DELETE FROM users WHERE id = ?;
-- -- name: SelectUserById :one
-- SELECT id, username FROM users WHERE id = ?;
-- -- name: SelectUserByName :one
-- SELECT id, username FROM users WHERE username = ?;
-- -- name: UpdateUserNameById :exec
-- UPDATE users SET username = ? WHERE id = ?;
