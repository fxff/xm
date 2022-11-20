CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA companies;

CREATE TYPE companies.COMPANY_TYPE AS ENUM (
    'corporation',
	'non-profit',
	'cooperative',
	'sole-proprietorship'
    );

CREATE TABLE companies.companies (
      id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY
    , "name" VARCHAR(15) NOT NULL
    , description VARCHAR
    , employees INT
    , registered BOOLEAN NOT NULL
    , type companies.COMPANY_TYPE
);

CREATE UNIQUE INDEX companies_unique_name ON companies.companies("name");