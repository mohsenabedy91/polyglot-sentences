-- Inserting data into access_controls
INSERT INTO access_controls (id, user_id, role_id)
VALUES (1, 1, 1),
       (2, 2, 2),
       (3, 3, 7),
       (4, 4, 3);

INSERT INTO access_controls (id, user_id, permission_id)
VALUES (5, 4, 1),
       (6, 4, 2),
       (7, 4, 3),
       (8, 2, 3);

SELECT setval('access_controls_id_seq', (SELECT MAX(id) FROM access_controls));