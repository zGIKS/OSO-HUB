-- Script para resetear las tablas de imágenes y empezar de cero
USE osohub;

-- Drop todas las tablas de imágenes
DROP TABLE IF EXISTS images_by_id;
DROP TABLE IF EXISTS images_by_date;
DROP TABLE IF EXISTS images_by_user;
DROP TABLE IF EXISTS image_counters;
DROP TABLE IF EXISTS likes_by_image;
DROP TABLE IF EXISTS reports_by_image;
DROP TABLE IF EXISTS users_by_id;