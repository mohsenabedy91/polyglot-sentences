-- Inserting data into roles
INSERT INTO roles (id, title, key, description, is_default, created_by, updated_by)
VALUES (1, 'Super Admin', 'SUPER_ADMIN', 'Super Admin role', TRUE, 1, 1),
       (2, 'Admin', 'ADMIN', 'Admin role', TRUE, 1, 1),
       (3, 'Manager', 'MANAGER', 'Manager role', TRUE, 1, 1),
       (4, 'Accountant', 'ACCOUNTANT', 'Accountant role', TRUE, 1, 1),
       (5, 'Supplier', 'SUPPLIER', 'Supplier role', TRUE, 1, 1),
       (6, 'Sales', 'SALES', 'Sales role', TRUE, 1, 1),
       (7, 'Staff', 'STAFF', 'Staff role', TRUE, 1, 1),
       (8, 'User', 'USER', 'User role', TRUE, 1, 1);

SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));