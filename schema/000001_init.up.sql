create table users (
 id            serial not null unique,
 name          varchar(255) not null,
 surname       varchar(255) not null,
 email         varchar(255) not null unique,
 password      varchar(255) not null,
 registered_at timestamp not null default now()
);

create table refresh_tokens(
 id            serial not null unique,
 user_id       int references users (id) on delete cascade not null,
 token         varchar(255) not null,
 expires_at    timestamp not null
);

create table files (
 id           serial not null unique,
 user_id      int references users (id) on delete cascade not null,
 name         varchar(255) not null,
 size         int not null,
 type         varchar(20) not null,
 content_type varchar(20) not null,
 url          varchar(255) not null,
 upload_at    timestamp not null default now()
);
