CREATE TABLE patients (
  id                 serial    not null primary key,
  name               varchar   not null,
  surname            varchar   not null,
  gender             varchar   not null,
  birth_year         integer   not null,
  phone              varchar   not null,
  email              varchar   not null unique,
  encrypted_password varchar   not null,
  check (birth_year between 1922 AND 2022),
  check (gender in ('Женский', 'Мужской'))
);

CREATE TABLE specializations (
  id     serial  not null primary key,
  name   varchar not null,
  salary integer not null,
  check (salary > 0)
);

CREATE TABLE diseases (
  id      serial not null primary key,
  name    varchar not null,
  spec_id integer   references specializations (id)
);

CREATE TABLE medicines (
  id   serial  not null primary key,
  name varchar not null,
  cost integer not null,
  check (cost > 0)
);

CREATE TABLE doctors (
  id                 serial    not null   primary key,
  name               varchar   not null,
  surname            varchar   not null,
  work_since         integer   not null,
  spec_id            integer   references specializations (id),
  email              varchar   not null unique,
  encrypted_password varchar   not null,
  check (work_since between 1942 AND 2022)
);

CREATE TABLE visits (
  id          serial  not null primary key,
  status      varchar not null,
  patient_id  integer references patients (id),
  doctor_id   integer references doctors (id),
  check (status in ('Done', 'Active'))
);

CREATE TABLE records (
  id          serial  not null primary key,
  visit_id    integer references visits (id),
  disease_id     integer references diseases (id),
  medicine_id    integer references medicines (id)
);

CREATE TABLE admins (
  id                 serial    not null primary key,
  name               varchar   not null,
  surname            varchar   not null,
  email              varchar   not null unique,
  encrypted_password varchar   not null
);

CREATE TABLE ops_stat
(
  operation char(1) not null,
  date timestamp not null,
  user_id text not null
);

--CREATE EXTENSION PLPYTHON3U;

CREATE OR REPLACE FUNCTION percent(integer)
RETURNS TABLE
(
    disease_id integer,
    percent    numeric
)
AS
$code$
    WITH cte AS 
    (
        SELECT r.disease_id, count(r.id)
        FROM records AS r JOIN diseases AS d ON r.disease_id = d.id
        WHERE d.spec_id = $1
        GROUP BY r.disease_id
    )
    SELECT disease_id, round((count::float / (select count(id) from records) * 100)::numeric, 2)
    FROM cte;
$code$
LANGUAGE SQL;

CREATE OR REPLACE FUNCTION logging()
RETURNS TRIGGER
AS
$code$
BEGIN
  IF (TG_OP = 'INSERT') THEN
    INSERT INTO ops_stat SELECT 'I', now(), user;
  ELSIF (TG_OP = 'DELETE') THEN
    INSERT INTO ops_stat SELECT 'S', now(), user;
  ELSIF (TG_OP = 'UPDATE') THEN
    INSERT INTO ops_stat SELECT 'S', now(), user;
  END IF;
  RETURN NULL;
END;
$code$
LANGUAGE PLPGSQL;

CREATE TRIGGER log AFTER INSERT OR DELETE OR UPDATE ON patients
FOR ROW EXECUTE FUNCTION logging();

CREATE TRIGGER log AFTER INSERT OR DELETE OR UPDATE ON doctors
FOR ROW EXECUTE FUNCTION logging();

CREATE TRIGGER log AFTER INSERT OR DELETE OR UPDATE ON admins
FOR ROW EXECUTE FUNCTION logging();

CREATE TRIGGER log AFTER INSERT OR DELETE OR UPDATE ON diseases
FOR ROW EXECUTE FUNCTION logging();

CREATE TRIGGER log AFTER INSERT OR DELETE OR UPDATE ON specializations
FOR ROW EXECUTE FUNCTION logging();

CREATE TRIGGER log AFTER INSERT OR DELETE OR UPDATE ON medicines
FOR ROW EXECUTE FUNCTION logging();

CREATE TRIGGER log AFTER INSERT OR DELETE OR UPDATE ON visits
FOR ROW EXECUTE FUNCTION logging();

CREATE TRIGGER log AFTER INSERT OR DELETE OR UPDATE ON records
FOR ROW EXECUTE FUNCTION logging();