-- public.person definition

-- Drop table

-- DROP TABLE public.person;

CREATE TABLE public.person (
	id  varchar NULL primary key,
	"name" varchar(200) NULL,
	gender varchar(200) NULL,
    email varchar NULL,
	date_create_at timestamptz NULL
);

CREATE TABLE public.webhook (
	id int4 NULL,
    email varchar NULL,
    url varchar(200) NULL,
	date_create_at timestamptz NULL
);


