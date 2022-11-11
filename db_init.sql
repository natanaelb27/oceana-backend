SELECT
    pg_terminate_backend(pid)
FROM
    pg_stat_activity
WHERE
    datname = 'oceana_db';

DROP DATABASE IF EXISTS oceana_db;

CREATE DATABASE oceana_db WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE = 'en_US.UTF-8';

\connect oceana_db;

CREATE SCHEMA public;

COMMENT ON SCHEMA public IS 'standard public schema';

SET default_tablespace = '';