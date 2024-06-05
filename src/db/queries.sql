-- table: rooms
-- name: InsertRoom :exec
INSERT INTO rooms (id, name) VALUES (?, ?);
-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id = ?;
-- name: SelectRoomById :one
SELECT id, name FROM rooms WHERE id = ?;
-- name: SelectRoomByName :one
SELECT id, name FROM rooms WHERE name = ?;
-- name: UpdateRoomNameById :exec
UPDATE rooms SET name = ? WHERE id = ?;

-- table: users
-- name: InsertUser :exec
INSERT INTO users (id, username) VALUES (?, ?);
-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
-- name: SelectUserById :one
SELECT id, username FROM users WHERE id = ?;
-- name: SelectUserByName :one
SELECT id, username FROM users WHERE username = ?;
-- name: UpdateUserNameById :exec
UPDATE users SET username = ? WHERE id = ?;
