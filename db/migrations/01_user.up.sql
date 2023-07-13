
create table user(
    id varchar(36) primary key not null, -- uuid
    name varchar(128) not null, -- username
    mail varchar(128) not null, -- email
    created bigint not null, -- creation time in `time.Now().Unix()`
    external_id varchar(64) not null, -- id on the external oauth provider
    type varchar(32) not null, -- GITHUB, DISCORD
    role varchar(16) not null
);