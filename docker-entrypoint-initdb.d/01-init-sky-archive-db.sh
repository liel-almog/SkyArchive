#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE SkyArchive;
EOSQL

# Switch to the 'SkyArchive' database and create tables
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "SkyArchive" <<-EOSQL
    -- Create sequences
    CREATE SEQUENCE IF NOT EXISTS files_file_id_seq;
    CREATE SEQUENCE IF NOT EXISTS users_id_seq;

    -- Create rooms table
CREATE TABLE IF NOT EXISTS public.files
(
    file_id integer NOT NULL DEFAULT nextval('files_file_id_seq'::regclass),
    display_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    original_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    uploaded_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    size bigint NOT NULL,
    mime_type character varying(150) COLLATE pg_catalog."default" NOT NULL,
    url text COLLATE pg_catalog."default",
    favorite boolean NOT NULL DEFAULT false,
    status character varying(50) COLLATE pg_catalog."default" NOT NULL DEFAULT 'PROCESSING'::character varying,
    user_id integer NOT NULL,
    CONSTRAINT files_pkey PRIMARY KEY (file_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id)
        REFERENCES public.users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.files
    OWNER to posgres;

    -- Create messages table
    CREATE TABLE IF NOT EXISTS public.users
(
    user_id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    username character varying(255) COLLATE pg_catalog."default" NOT NULL,
    email character varying(255) COLLATE pg_catalog."default" NOT NULL,
    password character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    CONSTRAINT users_email_key UNIQUE (email)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;

EOSQL