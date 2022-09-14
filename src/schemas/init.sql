INSERT INTO specializations (name, salary) VALUES 
('Терапевт',      15000),
('Хирург',        20000),
('Отоларинголог', 25000),
('Кардиолог',     30000),
('Онколог',       35000),
('Психиатр',      40000),
('Офтальмолог',   45000);

INSERT INTO diseases (name, spec_id) VALUES 
('ОРВИ',             1),
('ОРЗ',              1),
('Насморк',          1),

('Перелом',          2),
('Сотрясение мозга', 2),
('Микротравма',      2),

('Гайморит',         3),
('Бронхит',          3),
('Отит',             3),

('Ишемия',           4),
('Порок сердца',     4),
('Аритмия',          4),

('Рак легких',       5),
('Рак печени',       5),
('Рак мозга',        5),

('Депрессия',        6),
('Невроз',           6),
('Деменция',         6),

('Катаракта',        7),
('Глоукома',         7),
('Астигматизм',      7);

INSERT INTO medicines(name, cost) VALUES 
('Ношпа',                100),
('Аскорбиновая кислота', 200),
('Полисорб',             300),
('Аспирин',              400),
('Анальгин',             500),
('Мезим',                600),
('Лазолван',             700),
('Нурофен',              800),
('Ибупрофен',            900);

COPY patients(name, surname, gender, birth_year, phone, email, encrypted_password) 
FROM '/home/polina/course/http-rest-api/data/patients.txt' delimiter ',' csv NULL AS 'null';

COPY doctors(name, surname, work_since, spec_id, email, encrypted_password) 
FROM '/home/polina/course/http-rest-api/data/doctors.txt' delimiter ',' csv NULL AS 'null';

COPY admins(name, surname, email, encrypted_password) 
FROM '/home/polina/course/http-rest-api/data/admins.txt' delimiter ',' csv NULL AS 'null';

COPY visits(status, patient_id, doctor_id) 
FROM '/home/polina/course/http-rest-api/data/visits.txt' delimiter ',' csv NULL AS 'null';

COPY records(visit_id, disease_id, medicine_id) 
FROM '/home/polina/course/http-rest-api/data/records.txt' delimiter ',' csv NULL AS 'null';