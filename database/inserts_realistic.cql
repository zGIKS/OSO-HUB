-- 1. USUARIOS CON DATOS COMPLETAMENTE REALISTAS
INSERT INTO users_by_id (user_id, username, email, password_hash, profile_picture_url, bio, role, created_at)
VALUES (a1b2c3d4-e5f6-7890-1234-567890abcdef, 'mateo_dev', 'mateo.lopez@gmail.com', '$2b$12$LQv3c1yqBwrf2docu.ddaeb3UZzg7qUHXKkNrQePyOGz.aH5xKlTy', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face', 'Desarrollador Full Stack 💻 | Fotógrafo amateur 📸 | Amante del café ☕', 'user', '2025-01-15 10:30:00');

INSERT INTO users_by_id (user_id, username, email, password_hash, profile_picture_url, bio, role, created_at)
VALUES (b2c3d4e5-f6a7-8901-2345-678901bcdefb, 'ana_design', 'ana.martinez@outlook.com', '$2b$12$LQv3c1yqBwrf2docu.ddaeb3UZzg7qUHXKkNrQePyOGz.aH5xKlTy', 'https://images.unsplash.com/photo-1494790108755-2616b612b4b0?w=150&h=150&fit=crop&crop=face', 'UX/UI Designer 🎨 | Street photography enthusiast | Barcelona', 'user', '2025-02-20 14:20:00');

INSERT INTO users_by_id (user_id, username, email, password_hash, profile_picture_url, bio, role, created_at)
VALUES (c3d4e5f6-a7b8-9012-3456-789012cdefcd, 'carlos_photo', 'carlos.ruiz@yahoo.es', '$2b$12$LQv3c1yqBwrf2docu.ddaeb3UZzg7qUHXKkNrQePyOGz.aH5xKlTy', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face', 'Fotógrafo profesional 📷 | Madrid | Especialista en retratos', 'user', '2025-03-10 09:15:00');

-- 2. IMÁGENES CON user_profile_picture_url INCLUIDO
-- Imagen de mateo_dev
INSERT INTO images_by_id (image_id, day_bucket, uploaded_at, user_id, username, user_profile_picture_url, image_url, title)
VALUES (f1a2b3c4-d5e6-7890-1234-567890abcde1, '2025-06-30', '2025-06-30 18:30:00', a1b2c3d4-e5f6-7890-1234-567890abcdef, 'mateo_dev', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face', 'https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=800&h=600&fit=crop', 'Mi ritual matutino: café y código ☕💻');

INSERT INTO images_by_date (day_bucket, uploaded_at, image_id, user_id, username, user_profile_picture_url, image_url, title)
VALUES ('2025-06-30', '2025-06-30 18:30:00', f1a2b3c4-d5e6-7890-1234-567890abcde1, a1b2c3d4-e5f6-7890-1234-567890abcdef, 'mateo_dev', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face', 'https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=800&h=600&fit=crop', 'Mi ritual matutino: café y código ☕💻');

INSERT INTO images_by_user (user_id, uploaded_at, image_id, user_profile_picture_url, image_url, title)
VALUES (a1b2c3d4-e5f6-7890-1234-567890abcdef, '2025-06-30 18:30:00', f1a2b3c4-d5e6-7890-1234-567890abcde1, 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face', 'https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=800&h=600&fit=crop', 'Mi ritual matutino: café y código ☕💻');

-- Imagen de ana_design
INSERT INTO images_by_id (image_id, day_bucket, uploaded_at, user_id, username, user_profile_picture_url, image_url, title)
VALUES (f2a3b4c5-d6e7-8901-2345-678901bcdef2, '2025-06-30', '2025-06-30 19:15:00', b2c3d4e5-f6a7-8901-2345-678901bcdefb, 'ana_design', 'https://images.unsplash.com/photo-1494790108755-2616b612b4b0?w=150&h=150&fit=crop&crop=face', 'https://images.unsplash.com/photo-1449824913935-59a10b8d2000?w=800&h=600&fit=crop', 'Líneas perfectas en el centro de Barcelona 🏢');

INSERT INTO images_by_date (day_bucket, uploaded_at, image_id, user_id, username, user_profile_picture_url, image_url, title)
VALUES ('2025-06-30', '2025-06-30 19:15:00', f2a3b4c5-d6e7-8901-2345-678901bcdef2, b2c3d4e5-f6a7-8901-2345-678901bcdefb, 'ana_design', 'https://images.unsplash.com/photo-1494790108755-2616b612b4b0?w=150&h=150&fit=crop&crop=face', 'https://images.unsplash.com/photo-1449824913935-59a10b8d2000?w=800&h=600&fit=crop', 'Líneas perfectas en el centro de Barcelona 🏢');

INSERT INTO images_by_user (user_id, uploaded_at, image_id, user_profile_picture_url, image_url, title)
VALUES (b2c3d4e5-f6a7-8901-2345-678901bcdefb, '2025-06-30 19:15:00', f2a3b4c5-d6e7-8901-2345-678901bcdef2, 'https://images.unsplash.com/photo-1494790108755-2616b612b4b0?w=150&h=150&fit=crop&crop=face', 'https://images.unsplash.com/photo-1449824913935-59a10b8d2000?w=800&h=600&fit=crop', 'Líneas perfectas en el centro de Barcelona 🏢');

-- Imagen de carlos_photo
INSERT INTO images_by_id (image_id, day_bucket, uploaded_at, user_id, username, user_profile_picture_url, image_url, title)
VALUES (f3a4b5c6-d7e8-9012-3456-789012cdef34, '2025-06-30', '2025-06-30 20:47:00', c3d4e5f6-a7b8-9012-3456-789012cdefcd, 'carlos_photo', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face', 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800&h=600&fit=crop', 'Retrato urbano: luz natural en las calles de Madrid 📸');

INSERT INTO images_by_date (day_bucket, uploaded_at, image_id, user_id, username, user_profile_picture_url, image_url, title)
VALUES ('2025-06-30', '2025-06-30 20:47:00', f3a4b5c6-d7e8-9012-3456-789012cdef34, c3d4e5f6-a7b8-9012-3456-789012cdefcd, 'carlos_photo', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face', 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800&h=600&fit=crop', 'Retrato urbano: luz natural en las calles de Madrid 📸');

INSERT INTO images_by_user (user_id, uploaded_at, image_id, user_profile_picture_url, image_url, title)
VALUES (c3d4e5f6-a7b8-9012-3456-789012cdefcd, '2025-06-30 20:47:00', f3a4b5c6-d7e8-9012-3456-789012cdef34, 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face', 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800&h=600&fit=crop', 'Retrato urbano: luz natural en las calles de Madrid 📸');


--no tocar es para los nodos XD
ALTER KEYSPACE osohub WITH REPLICATION = {'class': 'NetworkTopologyStrategy', 'DC1': 3};