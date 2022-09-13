CREATE TABLE "users" (
	"id"	integer PRIMARY KEY AUTOINCREMENT,
	"created_at"	datetime,
	"updated_at"	datetime,
	"deleted_at"	datetime,
	"name"	varchar(255),
	"email"	varchar(255) UNIQUE,
	"hashed_password"	blob
);