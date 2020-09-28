CREATE TYPE gender AS ENUM
('male', 'female');

CREATE EXTENSION citext;

CREATE TABLE customers
(
	id SERIAL PRIMARY KEY,
	first_name VARCHAR(100) NOT NULL,
	last_name VARCHAR(100) NOT NULL,
	birth_date DATE NOT NULL 
	CHECK (birth_date < (CURRENT_DATE - interval '18' year)) 
	CHECK (birth_date > (CURRENT_DATE - interval '60' year)),
	gender VARCHAR NOT NULL,
	email citext UNIQUE NOT NULL,
	address VARCHAR
	(200)
);