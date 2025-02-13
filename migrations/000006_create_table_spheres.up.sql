CREATE TABLE IF NOT EXISTS executor_spheres (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  uniq VARCHAR(255),
  hash VARCHAR(255),
  parent_category INT REFERENCES executor_categories(id)
)
