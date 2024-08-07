CREATE TABLE IF NOT EXISTS rooms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    currently_playing INTEGER, -- Reference to the video table
    FOREIGN KEY (currently_playing) REFERENCES video(id)
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
);

CREATE TABLE IF NOT EXISTS video (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uri TEXT NOT NULL,
    local BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS room_history (
    room_id INTEGER,
    video_id INTEGER,
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (video_id) REFERENCES video(id)
);

CREATE TABLE IF NOT EXISTS user_history (
    user_id INTEGER,
    video_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (video_id) REFERENCES video(id)
);
