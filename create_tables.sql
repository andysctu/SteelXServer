CREATE TABLE game_history (gid int, start_time timestamp, end_time timestamp, game_type varchar(20), victor varchar(20), PRIMARY KEY(gid));
CREATE TABLE player_history (phid int, gid int, uid int, kills int, deaths int, assists int, team varchar(20), PRIMARY KEY(phid), FOREIGN KEY(gid) REFERENCES game_history(gid), FOREIGN KEY(uid) REFERENCES users(uid));
ALTER TABLE users ADD kills int, ADD deaths int, ADD assists int;
ALTER TABLE users ADD time_logged int;
