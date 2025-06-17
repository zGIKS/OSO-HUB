-- 1. USUARIOS CON DATOS COMPLETAMENTE REALISTAS
INSERT INTO users_by_id (user_id, username, email, password_hash, profile_picture_url, bio, role, created_at)
VALUES (a1b2c3d4-e5f6-7890-1234-567890abcdef, 'mateo_dev', 'mateo.lopez@gmail.com', '$2b$12$LQv3c1yqBwrf2docu.ddaeb3UZzg7qUHXKkNrQePyOGz.aH5xKlTy', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face', 'Desarrollador Full Stack 💻 | Fotógrafo amateur 📸 | Amante del café ☕', 'user', '2025-01-15 10:30:00');
INSERT INTO users_by_id (user_id, username, email, password_hash, profile_picture_url, bio, role, created_at)
VALUES (b2c3d4e5-f6a7-8901-2345-678901bcdefb, 'ana_design', 'ana.martinez@outlook.com', '$2b$12$LQv3c1yqBwrf2docu.ddaeb3UZzg7qUHXKkNrQePyOGz.aH5xKlTy', 'https://images.unsplash.com/photo-1494790108755-2616b612b4b0?w=150&h=150&fit=crop&crop=face', 'UX/UI Designer 🎨 | Street photography enthusiast | Barcelona', 'user', '2025-02-20 14:20:00');
INSERT INTO users_by_id (user_id, username, email, password_hash, profile_picture_url, bio, role, created_at)
VALUES (c3d4e5f6-a7b8-9012-3456-789012cdefac, 'carlos_nature', 'carlos.rivera@hotmail.com', '$2b$12$LQv3c1yqBwrf2docu.ddaeb3UZzg7qUHXKkNrQePyOGz.aH5xKlTy', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face', 'Fotografía de naturaleza 🌲 | Senderista los fines de semana | Madrid', 'user', '2025-03-10 09:15:00');
INSERT INTO users_by_id (user_id, username, email, password_hash, profile_picture_url, bio, role, created_at)
VALUES (d4e5f6a7-b8c9-0123-4567-890123defabd, 'laura_travel', 'laura.gomez@yahoo.com', '$2b$12$LQv3c1yqBwrf2docu.ddaeb3UZzg7qUHXKkNrQePyOGz.aH5xKlTy', 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150&h=150&fit=crop&crop=face', 'Travel blogger ✈️ | Explorando España una foto a la vez 🇪🇸', 'user', '2025-04-05 16:45:00');
INSERT INTO users_by_id (user_id, username, email, password_hash, profile_picture_url, bio, role, created_at)
VALUES (e5f6a7b8-c9d0-1234-5678-901234efabce, 'admin_system', 'admin@OSOHUBclone.com', '$2b$12$LQv3c1yqBwrf2docu.ddaeb3UZzg7qUHXKkNrQePyOGz.aH5xKlTy', 'https://images.unsplash.com/photo-1560250097-0b93528c311a?w=150&h=150&fit=crop&crop=face', 'Administrador del sistema OSOHUBClone', 'admin', '2025-01-01 00:00:00');

-- 2. IMÁGENES EN EL FEED GLOBAL (images_by_date) - HOY 2025-06-04
-- Café y trabajo matutino
INSERT INTO images_by_date (day_bucket, uploaded_at, image_id, user_id, username, image_url, title)
VALUES ('2025-06-04', '2025-06-04 08:30:00', f1a2b3c4-d5e6-7890-1234-567890abcde1, a1b2c3d4-e5f6-7890-1234-567890abcdef, 'mateo_dev', 'https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=800&h=600&fit=crop', 'Mi ritual matutino: café y código ☕💻');

-- Arquitectura urbana
INSERT INTO images_by_date (day_bucket, uploaded_at, image_id, user_id, username, image_url, title)
VALUES ('2025-06-04', '2025-06-04 10:15:00', f2a3b4c5-d6e7-8901-2345-678901bcdef2, b2c3d4e5-f6a7-8901-2345-678901bcdefb, 'ana_design', 'https://images.unsplash.com/photo-1449824913935-59a10b8d2000?w=800&h=600&fit=crop', 'Líneas perfectas en el centro de Barcelona 🏢');



-- 3. INSERTAR IMÁGENES EN PERFILES DE USUARIO (images_by_user)
-- Portfolio de Mateo
INSERT INTO images_by_user (user_id, uploaded_at, image_id, image_url, title)
VALUES (a1b2c3d4-e5f6-7890-1234-567890abcdef, '2025-06-04 08:30:00', f1a2b3c4-d5e6-7890-1234-567890abcde1, 'https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=800&h=600&fit=crop', 'Mi ritual matutino: café y código ☕💻');


-- Portfolio de Ana
INSERT INTO images_by_user (user_id, uploaded_at, image_id, image_url, title)
VALUES (b2c3d4e5-f6a7-8901-2345-678901bcdefb, '2025-06-04 10:15:00', f2a3b4c5-d6e7-8901-2345-678901bcdef2, 'https://images.unsplash.com/photo-1449824913935-59a10b8d2000?w=800&h=600&fit=crop', 'Líneas perfectas en el centro de Barcelona 🏢');



-- 4. INICIALIZAR CONTADORES DE LIKES REALISTAS
UPDATE image_counters SET likes = likes + 34, reports = reports + 1 WHERE image_id = f1a2b3c4-d5e6-7890-1234-567890abcde1;
UPDATE image_counters SET likes = likes + 67, reports = reports + 0 WHERE image_id = f2a3b4c5-d6e7-8901-2345-678901bcdef2;


-- 5. INSERTAR LIKES CRUZADOS REALISTAS (users liking each other's content)
-- Likes en foto de café de Mateo
INSERT INTO likes_by_image (image_id, user_id, liked_at)
VALUES (f1a2b3c4-d5e6-7890-1234-567890abcde1, b2c3d4e5-f6a7-8901-2345-678901bcdefb, '2025-06-04 08:45:00');

INSERT INTO likes_by_image (image_id, user_id, liked_at)
VALUES (f1a2b3c4-d5e6-7890-1234-567890abcde1, c3d4e5f6-a7b8-9012-3456-789012cdefac, '2025-06-04 09:15:00');

INSERT INTO likes_by_image (image_id, user_id, liked_at)
VALUES (f1a2b3c4-d5e6-7890-1234-567890abcde1, d4e5f6a7-b8c9-0123-4567-890123defabd, '2025-06-04 09:30:00');

-- Likes en arquitectura de Ana
INSERT INTO likes_by_image (image_id, user_id, liked_at)
VALUES (f2a3b4c5-d6e7-8901-2345-678901bcdef2, a1b2c3d4-e5f6-7890-1234-567890abcdef, '2025-06-04 10:30:00');

INSERT INTO likes_by_image (image_id, user_id, liked_at)
VALUES (f2a3b4c5-d6e7-8901-2345-678901bcdef2, c3d4e5f6-a7b8-9012-3456-789012cdefac, '2025-06-04 11:00:00');




-- 5 REPORTESS
INSERT INTO reports_by_image (image_id, report_id, reporter_id, reason, reported_at)
VALUES (f1a2b3c4-d5e6-7890-1234-567890abcde1, now(), e5f6a7b8-c9d0-1234-5678-901234efabce, 'Posible contenido duplicado', '2025-06-04 17:00:00');


-- Inserts para images_by_id (búsqueda directa por image_id)

INSERT INTO images_by_id (image_id, day_bucket, uploaded_at, user_id, username, image_url, title)
VALUES (f1a2b3c4-d5e6-7890-1234-567890abcde1, '2025-06-04', '2025-06-04 08:30:00', a1b2c3d4-e5f6-7890-1234-567890abcdef, 'mateo_dev', 'https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=800&h=600&fit=crop', 'Mi ritual matutino: café y código ☕💻');

INSERT INTO images_by_id (image_id, day_bucket, uploaded_at, user_id, username, image_url, title)
VALUES (f2a3b4c5-d6e7-8901-2345-678901bcdef2, '2025-06-04', '2025-06-04 10:15:00', b2c3d4e5-f6a7-8901-2345-678901bcdefb, 'ana_design', 'https://images.unsplash.com/photo-1449824913935-59a10b8d2000?w=800&h=600&fit=crop', 'Líneas perfectas en el centro de Barcelona 🏢');