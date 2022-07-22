CREATE TABLE books
(
    isbn   INTEGER,
    name   VARCHAR(50) NOT NULL,
    author VARCHAR(50) NOT NULL,
    CONSTRAINT books_pkey PRIMARY KEY (isbn)
)