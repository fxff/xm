CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP INDEX companies_unique_name;
DROP TABLE companies.companies;
DROP SCHEMA companies;

