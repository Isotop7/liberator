CREATE TABLE "books" (
	"id"	integer PRIMARY KEY AUTOINCREMENT,
	"created_at"	datetime,
	"updated_at"	datetime,
	"deleted_at"	datetime,
	"title"	varchar(255),
	"author"	varchar(255),
	"language"	varchar(255),
	"category"	varchar(255),
	"isbn10"	varchar(255),
	"isbn13"	varchar(255),
	"page_count"	integer,
	"rating"	integer,
	"review"	varchar(2048)
);