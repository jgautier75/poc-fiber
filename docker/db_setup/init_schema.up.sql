create database "poc-fiber" with encoding 'UTF-8' connection limit -1;
CREATE USER fiber_user WITH PASSWORD 'fiber_user';
GRANT ALL PRIVILEGES ON DATABASE "poc-fiber" TO fiber_user;
ALTER ROLE fiber_user NOSUPERUSER NOCREATEDB CREATEROLE INHERIT LOGIN;
grant all privileges on schema poc-fiber.public to fiber_user;
