-- Inserting data into roles
INSERT INTO roles (id, title, key, description, created_by, updated_by)
VALUES (1, 'Super Admin', 'SUPER_ADMIN', 'Super Admin role', 1, 1),
       (2, 'Admin', 'ADMIN', 'Admin role', 1, 1),
       (3, 'Manager', 'MANAGER', 'Manager role', 1, 1),
       (4, 'Accountant', 'ACCOUNTANT', 'Accountant role', 1, 1),
       (5, 'Supplier', 'SUPPLIER', 'Supplier role', 1, 1),
       (6, 'Sales', 'SALES', 'Sales role', 1, 1),
       (7, 'Staff', 'STAFF', 'Staff role', 1, 1),
       (8, 'User', 'USER', 'User role', 1, 1);

SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));