-- Inserting data into permissions
INSERT INTO permissions (id, title, key, "group", description, created_by, updated_by)
VALUES (1, 'Create user', 'CREATE_USER', 'user', 'Create a new user', 1, 1),
       (2, 'Read user', 'READ_USER', 'user', 'Read user information', 1, 1),
       (3, 'Update user', 'UPDATE_USER', 'user', 'Update user information', 1, 1),
       (4, 'Delete user', 'DELETE_USER', 'user', 'Delete user', 1, 1),
       (5, 'Create role', 'CREATE_ROLE', 'role', 'Create a new role', 1, 1),
       (6, 'Read role', 'READ_ROLE', 'role', 'Read role information', 1, 1),
       (7, 'Update role', 'UPDATE_ROLE', 'role', 'Update role information', 1, 1),
       (8, 'Delete role', 'DELETE_ROLE', 'role', 'Delete role', 1, 1),
       (9, 'Create permission', 'CREATE_PERMISSION', 'permission', 'Create a new permission', 1, 1),
       (10, 'Read permission', 'READ_PERMISSION', 'permission', 'Read permission information', 1, 1),
       (11, 'Update permission', 'UPDATE_PERMISSION', 'permission', 'Update permission information', 1, 1),
       (12, 'Delete permission', 'DELETE_PERMISSION', 'permission', 'Delete permission', 1, 1),
       (13, 'Assign Roles To User', 'ASSIGN_ROLES_TO_USER', 'access_control', 'Assign roles to user', 1, 1),
       (14, 'Read User Roles', 'READ_USER_ROLES', 'access_control', 'Read user roles information', 1, 1),
       (15, 'Assign Permissions To Role', 'ASSIGN_PERMISSIONS_TO_ROLE', 'access_control', 'Assign permissions to role ', 1, 1),
       (16, 'Read Role Permissions', 'READ_ROLE_PERMISSIONS', 'access_control', 'Read role permissions information', 1, 1);

SELECT setval('permissions_id_seq', (SELECT MAX(id) FROM permissions));