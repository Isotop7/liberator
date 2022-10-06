CREATE TABLE "books_users_assignment" (
    "id" integer PRIMARY KEY AUTOINCREMENT,
    "created_at"	datetime,
	"updated_at"	datetime,
	"deleted_at"	datetime,
    "user_id" integer REFERENCES users(id),
    "books_id" integer REFERENCES books(id),
    "status" varchar(16),
    "pages_read" integer,
    "rating" integer
);