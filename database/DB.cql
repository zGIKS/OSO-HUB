
CREATE KEYSPACE osohub WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE osohub;



-- 7. Tabla para búsqueda directa por image_id
CREATE TABLE IF NOT EXISTS images_by_id (
  image_id uuid PRIMARY KEY,
  day_bucket text,
  uploaded_at timestamp,
  user_id uuid,
  username text,
  image_url text,
  title text
);
-- 1. Tabla de usuarios
CREATE TABLE IF NOT EXISTS users_by_id (
  user_id uuid PRIMARY KEY,
  username text,
  email text,
  password_hash text,
  profile_picture_url text,
  bio text,
  role text, -- 'user' o 'admin'
  created_at timestamp
);

-- 2. Tabla para feed global (home) con date bucketing
CREATE TABLE IF NOT EXISTS images_by_date (
  day_bucket text, -- formato 'YYYY-MM-DD'
  uploaded_at timestamp,
  image_id uuid,
  user_id uuid,
  username text,
  image_url text,
  title text,
  PRIMARY KEY ((day_bucket), uploaded_at, image_id)
) WITH CLUSTERING ORDER BY (uploaded_at DESC, image_id ASC);

-- Índice secundario para email en users_by_id
CREATE INDEX IF NOT EXISTS users_by_id_email_idx ON users_by_id (email);

-- 3. Tabla para publicaciones por usuario (perfil)
CREATE TABLE IF NOT EXISTS images_by_user (
  user_id uuid,
  uploaded_at timestamp,
  image_id uuid,
  image_url text,
  title text,
  PRIMARY KEY (user_id, uploaded_at, image_id)
) WITH CLUSTERING ORDER BY (uploaded_at DESC, image_id ASC);

-- 4. Tabla de contadores por imagen (likes y reportes)
CREATE TABLE IF NOT EXISTS image_counters (
  image_id uuid PRIMARY KEY,
  likes counter,
  reports counter
);

-- 5. Tabla de reportes por imagen (detalle)
CREATE TABLE IF NOT EXISTS reports_by_image (
  image_id uuid,
  report_id timeuuid,
  reporter_id uuid,
  reason text,
  reported_at timestamp,
  PRIMARY KEY (image_id, report_id)
) WITH CLUSTERING ORDER BY (report_id DESC);

-- 6. Tabla de likes por imagen (para evitar múltiples likes por usuario)
CREATE TABLE IF NOT EXISTS likes_by_image (
  image_id uuid,
  user_id uuid,
  liked_at timestamp,
  PRIMARY KEY (image_id, user_id)
);


