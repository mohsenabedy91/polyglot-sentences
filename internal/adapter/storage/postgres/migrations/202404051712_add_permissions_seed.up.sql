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
       (9, 'Read permission', 'READ_PERMISSION', 'permission', 'Read permission information', 1, 1),
       (10, 'Sync Roles With User', 'SYNC_ROLES_WITH_USER', 'access_control', 'Sync roles With user', 1, 1),
       (11, 'Read User Roles', 'READ_USER_ROLES', 'access_control', 'Read user roles information', 1, 1),
       (12, 'Sync Permissions With Role', 'SYNC_PERMISSIONS_WITH_ROLE', 'access_control', 'Sync permissions With role ', 1, 1),
       (13, 'Read Role Permissions', 'READ_ROLE_PERMISSIONS', 'access_control', 'Read role permissions information', 1, 1);

SELECT setval('permissions_id_seq', (SELECT MAX(id) FROM permissions));