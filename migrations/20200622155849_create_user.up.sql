CREATE TABLE users (
  id bigserial not null primary key,
  login_name varchar not null unique,
  encrypted_password varchar not null,
  jobcode varchar not null,
  email varchar not null,
  phone varchar not null,
  user_name varchar not null
);
