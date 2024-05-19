-- Inserting data into the users table
INSERT INTO users (id, language_id, first_name, last_name, email, password, status, created_by, updated_by)
VALUES (1, 1, 'John', 'Doe', 'john.doe@gmail.com', '$2a$11$UgeDhqEyJ8ychQqApHlNleQ4pzs0nh3wSy024SHaTDI7DLOiUeUlu',
        'ACTIVE', 1, 1), -- password: QWer123!@#
       (2, 1, 'Jane', 'Doe', 'jane.doe@gmail.com', '$2a$11$UgeDhqEyJ8ychQqApHlNleQ4pzs0nh3wSy024SHaTDI7DLOiUeUlu',
        'ACTIVE', 1, 1), -- password: QWer123!@#
       (3, 1, 'Alice', 'Smith', 'alice@gmail.com', '$2a$11$UgeDhqEyJ8ychQqApHlNleQ4pzs0nh3wSy024SHaTDI7DLOiUeUlu',
        'ACTIVE', 1, 1), -- password: QWer123!@#
       (4, 1, 'Bob', 'Smith', 'b_smith@gmail.com', '$2a$11$UgeDhqEyJ8ychQqApHlNleQ4pzs0nh3wSy024SHaTDI7DLOiUeUlu',
        'ACTIVE', 1, 1); -- password: QWer123!@#

SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));