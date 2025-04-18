INSERT INTO public.users (id, created_at, updated_at, deleted_at, "name", email, "password", "role")
VALUES (
  1,
  '2025-04-16 11:57:55.405',
  '2025-04-16 11:57:55.405',
  NULL,
  'admin',
  'admin@admin.com',
  '$2a$10$v2M/FAuMLsP9spkQDYi8IeJDlluI58vJyU.jwTMaNIe0k1GALyIc2',
  'admin'
)
ON CONFLICT (email) DO NOTHING;
