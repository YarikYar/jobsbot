CREATE TABLE IF NOT EXISTS offer_categories (
  id serial primary key,
  name varchar(255) not null,
  data varchar(255) not null,
  hash varchar(255) not null
);
