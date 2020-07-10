CREATE TABLE users (
  id bigserial not null primary key,
  encrypted_password varchar not null,
  login_name varchar not null unique,
  jobcode varchar not null,
  email varchar not null,
  phone varchar not null,
  user_name varchar not null,
  uuid varchar not null,
  roles varchar not null
);

CREATE TABLE apps (
  id bigserial not null primary key,
  app_name varchar not null,
  app_system_name varchar not null unique,
  app_id varchar not null unique,
  app_state boolean not null,
  rating int not null,
  app_category varchar not null,
  token varchar not null
);

CREATE TABLE launch (
  id bigserial not null primary key,
  app_system_name varchar not null unique,
  app_id varchar not null unique,
  pid int not null
);
