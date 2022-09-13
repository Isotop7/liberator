CREATE TABLE "sessions" (
	"token"	TEXT,
	"data"	BLOB NOT NULL,
	"expiry"	REAL NOT NULL,
	PRIMARY KEY("token")
);

CREATE INDEX idx_sessions_expiry 
ON sessions(expiry); 