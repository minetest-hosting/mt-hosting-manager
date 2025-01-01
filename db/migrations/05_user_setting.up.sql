-- a per-user setting
create table user_setting(
    user_id varchar(36) not null references public.user(id) on delete cascade,
    key varchar(64) not null default '',
    value varchar(512) not null default '',
    primary key (user_id, key)
);