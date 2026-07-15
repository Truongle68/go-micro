-- Seed data for user-service
-- Admin user:
--   username: admin
--   phone/identifier: +84900000001
--   password: password123
--   full_name: Admin User

-- Normal user:
--   username: user1
--   phone/identifier: +84900000002
--   password: password123
--   full_name: Normal User

INSERT INTO users (id, username, status, role) VALUES 
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin', 'verified', 'admin'),
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'user1', 'verified', 'customer')
ON CONFLICT (id) DO NOTHING;

INSERT INTO user_credentials (user_id, type, identifier, secret_hash, is_verified, is_primary) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'phone', '+84900000001', '$2a$10$K3HgKiI5hmRF4PdDsza92.uPImIQInddqD4WQORaVVqdEGUiEEjwq', true, true),
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'phone', '+84900000002', '$2a$10$K3HgKiI5hmRF4PdDsza92.uPImIQInddqD4WQORaVVqdEGUiEEjwq', true, true)
ON CONFLICT (identifier) DO NOTHING;

INSERT INTO profiles (user_id, full_name, avatar_url, gender, dob) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Admin User', '', 'male', '1990-01-01'),
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Normal User', '', 'female', '1995-05-05')
ON CONFLICT (user_id) DO NOTHING;
