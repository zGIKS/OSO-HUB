USE OSOHUB;

CREATE KEYSPACE IF NOT EXISTS OSOHUB
  WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
  };

CREATE TABLE IF NOT EXISTS users_by_id (
  user_id       uuid PRIMARY KEY,
  username      text,
  email         text,
  password_hash text,
  full_name     text,
  bio           text,
  avatar_url    text,
  created_at    timestamp
);

CREATE TABLE IF NOT EXISTS followers_by_user (
  user_id     uuid,
  follower_id uuid,
  followed_at timestamp,
  PRIMARY KEY (user_id, follower_id)
);

CREATE TABLE IF NOT EXISTS followees_by_user (
  user_id     uuid,
  followee_id uuid,
  followed_at timestamp,
  PRIMARY KEY (user_id, followee_id)
);

CREATE TABLE IF NOT EXISTS posts_by_user (
  user_id      uuid,
  created_at   timestamp,
  post_id      uuid,
  title        text,
  description  text,
  image_urls   list<text>,
  PRIMARY KEY (user_id, created_at, post_id)
) WITH CLUSTERING ORDER BY (created_at DESC, post_id ASC);

CREATE TABLE IF NOT EXISTS posts_by_date_bucket (
  date_bucket  text,
  created_at   timestamp,
  post_id      uuid,
  user_id      uuid,
  title        text,
  description  text,
  image_urls   list<text>,
  PRIMARY KEY (date_bucket, created_at, post_id)
) WITH CLUSTERING ORDER BY (created_at DESC, post_id ASC);

CREATE TABLE IF NOT EXISTS posts_by_keyword (
  keyword      text,
  created_at   timestamp,
  post_id      uuid,
  user_id      uuid,
  PRIMARY KEY (keyword, created_at, post_id)
) WITH CLUSTERING ORDER BY (created_at DESC, post_id ASC);

CREATE TABLE IF NOT EXISTS comments_by_post (
  post_id      uuid,
  created_at   timestamp,
  comment_id   timeuuid,
  commenter_id uuid,
  content      text,
  PRIMARY KEY (post_id, created_at, comment_id)
) WITH CLUSTERING ORDER BY (created_at DESC, comment_id DESC);

CREATE TABLE IF NOT EXISTS likes_count (
  post_id uuid PRIMARY KEY,
  likes   counter
);

CREATE TABLE IF NOT EXISTS likes_by_post (
  post_id  uuid,
  user_id  uuid,
  liked_at timestamp,
  PRIMARY KEY (post_id, user_id)
);

CREATE TABLE IF NOT EXISTS feed_by_user (
  user_id          uuid,
  post_created_at  timestamp,
  post_id          uuid,
  author_id        uuid,
  title            text,
  PRIMARY KEY (user_id, post_created_at, post_id)
) WITH CLUSTERING ORDER BY (post_created_at DESC, post_id ASC);
