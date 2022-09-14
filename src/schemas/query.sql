create role patient with login password '111';
grant select, insert on table patients to patient;
grant select,insert on table visits to patient;
grant select on table records to patient;

create role doctor with login password '222';
grant select, insert on table doctors to doctor;
grant select on table medicines to doctor;
grant select on table diseases to doctor;
grant select, update on table visits to doctor;
grant select, insert on table records to doctor;

create role admin with login password '333';
grant select, insert on table admins to admin;
grant select on table patient to admin;
grant select on table medicines to doctor;
grant select on table diseases to doctor;
grant select, insert on table doctors to admin;
grant select, insert on table records to admin;