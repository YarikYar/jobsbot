CREATE TABLE IF NOT EXISTS posts (
  id SERIAL PRIMARY KEY,
  post_type TEXT NOT NULL,
  title TEXT NOT NULL,
 post_link TEXT NOT NULL,
 post_id INT NOT NULL,
  user_id BIGINT NOT NULL,
 created_at timestamp with time zone

)
