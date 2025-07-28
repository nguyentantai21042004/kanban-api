-- Initialize basic data for Kanban application
-- This file contains sample data for users, boards, lists, cards, labels, and activities

-- Sample Users
INSERT INTO users (id, username, full_name, password_hash, avatar_url, is_active) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'admin', 'System Administrator', '$2a$10$hashedpassword123', 'https://example.com/avatars/admin.jpg', true),
    ('550e8400-e29b-41d4-a716-446655440002', 'john.doe', 'John Doe', '$2a$10$hashedpassword456', 'https://example.com/avatars/john.jpg', true),
    ('550e8400-e29b-41d4-a716-446655440003', 'jane.smith', 'Jane Smith', '$2a$10$hashedpassword789', 'https://example.com/avatars/jane.jpg', true),
    ('550e8400-e29b-41d4-a716-446655440004', 'mike.wilson', 'Mike Wilson', '$2a$10$hashedpassword012', 'https://example.com/avatars/mike.jpg', true),
    ('550e8400-e29b-41d4-a716-446655440005', 'sarah.jones', 'Sarah Jones', '$2a$10$hashedpassword345', 'https://example.com/avatars/sarah.jpg', true)
ON CONFLICT (username) DO NOTHING;

-- Sample Boards
INSERT INTO boards (id, title, description) VALUES
    ('660e8400-e29b-41d4-a716-446655440001', 'Product Development', 'Main board for product development tasks and features'),
    ('660e8400-e29b-41d4-a716-446655440002', 'Marketing Campaign', 'Board for managing marketing campaigns and content'),
    ('660e8400-e29b-41d4-a716-446655440003', 'Bug Tracking', 'Board for tracking and resolving software bugs'),
    ('660e8400-e29b-41d4-a716-446655440004', 'Customer Support', 'Board for managing customer support tickets')
ON CONFLICT DO NOTHING;

-- Sample Labels for Product Development Board
INSERT INTO labels (id, board_id, name, color) VALUES
    ('770e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001', 'Frontend', '#FF6B6B'),
    ('770e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440001', 'Backend', '#4ECDC4'),
    ('770e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440001', 'Database', '#45B7D1'),
    ('770e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440001', 'UI/UX', '#96CEB4'),
    ('770e8400-e29b-41d4-a716-446655440005', '660e8400-e29b-41d4-a716-446655440001', 'Testing', '#FFEAA7')
ON CONFLICT DO NOTHING;

-- Sample Labels for Marketing Campaign Board
INSERT INTO labels (id, board_id, name, color) VALUES
    ('770e8400-e29b-41d4-a716-446655440006', '660e8400-e29b-41d4-a716-446655440002', 'Social Media', '#FF8A80'),
    ('770e8400-e29b-41d4-a716-446655440007', '660e8400-e29b-41d4-a716-446655440002', 'Email', '#80D8FF'),
    ('770e8400-e29b-41d4-a716-446655440008', '660e8400-e29b-41d4-a716-446655440002', 'Content', '#FFD180'),
    ('770e8400-e29b-41d4-a716-446655440009', '660e8400-e29b-41d4-a716-446655440002', 'Analytics', '#A7FFEB')
ON CONFLICT DO NOTHING;

-- Sample Lists for Product Development Board
INSERT INTO lists (id, board_id, title, position, is_archived) VALUES
    ('880e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001', 'Backlog', 1.0, false),
    ('880e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440001', 'To Do', 2.0, false),
    ('880e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440001', 'In Progress', 3.0, false),
    ('880e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440001', 'Review', 4.0, false),
    ('880e8400-e29b-41d4-a716-446655440005', '660e8400-e29b-41d4-a716-446655440001', 'Done', 5.0, false)
ON CONFLICT DO NOTHING;

-- Sample Lists for Marketing Campaign Board
INSERT INTO lists (id, board_id, title, position, is_archived) VALUES
    ('880e8400-e29b-41d4-a716-446655440006', '660e8400-e29b-41d4-a716-446655440002', 'Ideas', 1.0, false),
    ('880e8400-e29b-41d4-a716-446655440007', '660e8400-e29b-41d4-a716-446655440002', 'Planning', 2.0, false),
    ('880e8400-e29b-41d4-a716-446655440008', '660e8400-e29b-41d4-a716-446655440002', 'In Progress', 3.0, false),
    ('880e8400-e29b-41d4-a716-446655440009', '660e8400-e29b-41d4-a716-446655440002', 'Published', 4.0, false)
ON CONFLICT DO NOTHING;

-- Sample Cards for Product Development Board
INSERT INTO cards (id, list_id, title, description, position, due_date, priority, labels, is_archived) VALUES
    -- Backlog Cards
    ('990e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440001', 'Implement User Authentication', 'Add JWT-based authentication system with refresh tokens', 1.0, '2024-02-15 17:00:00+07', 'high', '[{"id": "770e8400-e29b-41d4-a716-446655440002", "name": "Backend", "color": "#4ECDC4"}]', false),
    ('990e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440001', 'Design Mobile App UI', 'Create wireframes and mockups for mobile application', 2.0, '2024-02-20 17:00:00+07', 'medium', '[{"id": "770e8400-e29b-41d4-a716-446655440004", "name": "UI/UX", "color": "#96CEB4"}]', false),
    
    -- To Do Cards
    ('990e8400-e29b-41d4-a716-446655440003', '880e8400-e29b-41d4-a716-446655440002', 'Setup Database Schema', 'Create database tables and relationships for the application', 1.0, '2024-02-10 17:00:00+07', 'high', '[{"id": "770e8400-e29b-41d4-a716-446655440003", "name": "Database", "color": "#45B7D1"}]', false),
    ('990e8400-e29b-41d4-a716-446655440004', '880e8400-e29b-41d4-a716-446655440002', 'Create Login Page', 'Implement responsive login page with form validation', 2.0, '2024-02-12 17:00:00+07', 'medium', '[{"id": "770e8400-e29b-41d4-a716-446655440001", "name": "Frontend", "color": "#FF6B6B"}]', false),
    
    -- In Progress Cards
    ('990e8400-e29b-41d4-a716-446655440005', '880e8400-e29b-41d4-a716-446655440003', 'API Development', 'Building RESTful APIs for user management', 1.0, '2024-02-08 17:00:00+07', 'high', '[{"id": "770e8400-e29b-41d4-a716-446655440002", "name": "Backend", "color": "#4ECDC4"}]', false),
    ('990e8400-e29b-41d4-a716-446655440006', '880e8400-e29b-41d4-a716-446655440003', 'Unit Testing', 'Write unit tests for core functionality', 2.0, '2024-02-14 17:00:00+07', 'medium', '[{"id": "770e8400-e29b-41d4-a716-446655440005", "name": "Testing", "color": "#FFEAA7"}]', false),
    
    -- Review Cards
    ('990e8400-e29b-41d4-a716-446655440007', '880e8400-e29b-41d4-a716-446655440004', 'Code Review: User Module', 'Review the user management module implementation', 1.0, '2024-02-06 17:00:00+07', 'medium', '[{"id": "770e8400-e29b-41d4-a716-446655440002", "name": "Backend", "color": "#4ECDC4"}]', false),
    
    -- Done Cards
    ('990e8400-e29b-41d4-a716-446655440008', '880e8400-e29b-41d4-a716-446655440005', 'Project Setup', 'Initial project setup with basic configuration', 1.0, '2024-02-01 17:00:00+07', 'low', '[]', false),
    ('990e8400-e29b-41d4-a716-446655440009', '880e8400-e29b-41d4-a716-446655440005', 'Requirements Gathering', 'Completed requirements analysis and documentation', 2.0, '2024-01-30 17:00:00+07', 'medium', '[]', false)
ON CONFLICT DO NOTHING;

-- Sample Cards for Marketing Campaign Board
INSERT INTO cards (id, list_id, title, description, position, due_date, priority, labels, is_archived) VALUES
    -- Ideas Cards
    ('990e8400-e29b-41d4-a716-446655440010', '880e8400-e29b-41d4-a716-446655440006', 'Social Media Campaign', 'Launch campaign across Facebook, Instagram, and LinkedIn', 1.0, '2024-02-25 17:00:00+07', 'high', '[{"id": "770e8400-e29b-41d4-a716-446655440006", "name": "Social Media", "color": "#FF8A80"}]', false),
    
    -- Planning Cards
    ('990e8400-e29b-41d4-a716-446655440011', '880e8400-e29b-41d4-a716-446655440007', 'Email Newsletter Design', 'Design monthly newsletter template', 1.0, '2024-02-18 17:00:00+07', 'medium', '[{"id": "770e8400-e29b-41d4-a716-446655440007", "name": "Email", "color": "#80D8FF"}]', false),
    
    -- In Progress Cards
    ('990e8400-e29b-41d4-a716-446655440012', '880e8400-e29b-41d4-a716-446655440008', 'Blog Content Creation', 'Writing blog posts about product features', 1.0, '2024-02-16 17:00:00+07', 'medium', '[{"id": "770e8400-e29b-41d4-a716-446655440008", "name": "Content", "color": "#FFD180"}]', false),
    
    -- Published Cards
    ('990e8400-e29b-41d4-a716-446655440013', '880e8400-e29b-41d4-a716-446655440009', 'Website Banner Update', 'Updated website banner with new promotion', 1.0, '2024-02-05 17:00:00+07', 'low', '[{"id": "770e8400-e29b-41d4-a716-446655440008", "name": "Content", "color": "#FFD180"}]', false)
ON CONFLICT DO NOTHING;

-- Sample Card Activities
INSERT INTO card_activities (id, card_id, action_type, old_data, new_data) VALUES
    ('aa0e8400-e29b-41d4-a716-446655440001', '990e8400-e29b-41d4-a716-446655440001', 'created', NULL, '{"title": "Implement User Authentication", "list_id": "880e8400-e29b-41d4-a716-446655440001"}'),
    ('aa0e8400-e29b-41d4-a716-446655440002', '990e8400-e29b-41d4-a716-446655440003', 'created', NULL, '{"title": "Setup Database Schema", "list_id": "880e8400-e29b-41d4-a716-446655440002"}'),
    ('aa0e8400-e29b-41d4-a716-446655440003', '990e8400-e29b-41d4-a716-446655440005', 'created', NULL, '{"title": "API Development", "list_id": "880e8400-e29b-41d4-a716-446655440003"}'),
    ('aa0e8400-e29b-41d4-a716-446655440004', '990e8400-e29b-41d4-a716-446655440005', 'moved', '{"list_id": "880e8400-e29b-41d4-a716-446655440002"}', '{"list_id": "880e8400-e29b-41d4-a716-446655440003"}'),
    ('aa0e8400-e29b-41d4-a716-446655440005', '990e8400-e29b-41d4-a716-446655440007', 'updated', '{"priority": "low"}', '{"priority": "medium"}'),
    ('aa0e8400-e29b-41d4-a716-446655440006', '990e8400-e29b-41d4-a716-446655440008', 'moved', '{"list_id": "880e8400-e29b-41d4-a716-446655440004"}', '{"list_id": "880e8400-e29b-41d4-a716-446655440005"}'),
    ('aa0e8400-e29b-41d4-a716-446655440007', '990e8400-e29b-41d4-a716-446655440010', 'created', NULL, '{"title": "Social Media Campaign", "list_id": "880e8400-e29b-41d4-a716-446655440006"}'),
    ('aa0e8400-e29b-41d4-a716-446655440008', '990e8400-e29b-41d4-a716-446655440013', 'moved', '{"list_id": "880e8400-e29b-41d4-a716-446655440008"}', '{"list_id": "880e8400-e29b-41d4-a716-446655440009"}')
ON CONFLICT DO NOTHING;

-- Add some comments to activities
INSERT INTO card_activities (id, card_id, action_type, old_data, new_data) VALUES
    ('aa0e8400-e29b-41d4-a716-446655440009', '990e8400-e29b-41d4-a716-446655440001', 'commented', NULL, '{"comment": "This is a high priority task that needs to be completed before the next release."}'),
    ('aa0e8400-e29b-41d4-a716-446655440010', '990e8400-e29b-41d4-a716-446655440005', 'commented', NULL, '{"comment": "API endpoints are working well. Ready for testing phase."}'),
    ('aa0e8400-e29b-41d4-a716-446655440011', '990e8400-e29b-41d4-a716-446655440010', 'commented', NULL, '{"comment": "Great idea! Let''s start with Facebook and Instagram first."}')
ON CONFLICT DO NOTHING; 