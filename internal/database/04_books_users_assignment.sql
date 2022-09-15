CREATE TABLE "books_users_assignment" (
    "id" integer PRIMARY KEY AUTOINCREMENT,
    "user_id" integer REFERENCES users(id),
    "books_id" integer REFERENCES books(id),
    "status" integer,
    "pages_read" integer,
    "rating" integer
);