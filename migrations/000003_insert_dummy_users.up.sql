INSERT INTO users (id, name, birthday, email, password_hash) VALUES
    ('00000000-0000-0000-0000-000000000000', 'John Doe', NOW()-INTERVAL '30 years'-INTERVAL '15 days', 'fHkIj@example.com', '$2a$10$KUhxHLsyTIYFBPhuObrsl.GE0tScBD9ix6LsZ12Yxpwp8lPfHJGJ6'),
    ('00000000-0000-0000-0000-000000000001', 'Jane Doe', NOW()-INTERVAL '15 years'-INTERVAL '5 days', 'qg6Z8@example.com', '$2a$10$KUhxHLsyTIYFBPhuObrsl.GE0tScBD9ix6LsZ12Yxpwp8lPfHJGJ6'),
    ('00000000-0000-0000-0000-000000000002', 'John Smith', NOW()-INTERVAL '45 years'-INTERVAL '10 days', 'HjZVz@example.com', '$2a$10$KUhxHLsyTIYFBPhuObrsl.GE0tScBD9ix6LsZ12Yxpwp8lPfHJGJ6'),
    ('00000000-0000-0000-0000-000000000003', 'Jane Smith', NOW()-INTERVAL '50 years', '5hBp1@example.com', '$2a$10$KUhxHLsyTIYFBPhuObrsl.GE0tScBD9ix6LsZ12Yxpwp8lPfHJGJ6')