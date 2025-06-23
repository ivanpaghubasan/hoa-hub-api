-- ROLES
INSERT INTO roles (name, description) VALUES
('admin', 'Administrator with full permissions'),
('member', 'Regular homeowner member'),
('treasurer', 'Manages financial records');

-- PERMISSIONS
INSERT INTO permissions (name, description) VALUES
('manage_users', 'Add, edit, delete users'),
('manage_properties', 'Add, edit, delete properties'),
('view_reports', 'Access reports and analytics'),
('manage_finances', 'Handle financial data');

-- admin gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'admin';

-- member gets view_reports
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'member' AND p.name = 'view_reports';

-- treasurer gets view_reports + manage_finances
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'treasurer' AND p.name IN ('view_reports', 'manage_finances');