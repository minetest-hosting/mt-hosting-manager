PRAGMA foreign_keys = ON;

-- registered user
create table user(
    id varchar(36) primary key not null, -- uuid
    state varchar(16) not null default 'ACTIVE',
    name varchar(128) not null, -- username
    mail varchar(128) not null, -- email
    created bigint not null, -- creation time in `time.Now().Unix()`
    external_id varchar(64) not null, -- id on the external oauth provider
    type varchar(32) not null, -- GITHUB, DISCORD
    role varchar(16) not null, -- ADMIN / USER
    credits bigint not null default 0, -- credits in millicents USD
    max_credits bigint not null default 100, -- max credits in millicents USD
    warn_credits bigint not null default 5 -- warn if the credits fall below this threshold
);

-- available node types to select
create table node_type(
    id varchar(36) primary key not null, -- uuid
    state varchar(16) not null default 'INACTIVE',
    order_id int not null default 1, -- order id
    provider varchar(16) not null, -- HETZNER
    server_type varchar(16) not null, -- provider server-type: cx11
    name varchar(128) not null default '', -- name of the node
    description varchar(1024) not null default '', -- description of the node
    cost_per_hour bigint not null default 0.01, -- cost in millicents USD per hour
    max_recommended_instances int not null default 2, -- max number of recommended minetest instances on this host
    max_instances int not null default 4 -- max number of allowed minetest instances on this host
);

-- a node set up by a user
create table user_node(
    id varchar(36) primary key not null, -- uuid
    user_id varchar(36) not null references user(id),
    node_type_id varchar(36) not null references node_type(id),
    created bigint not null, -- creation time in `time.Now().Unix()`
    state varchar(16) not null default 'CREATED',
    name varchar(64) not null, -- name of the host, used for dns registration (A, AAAA record)
    ipv4 varchar(32) not null,
    ipv6 varchar(128) not null
);

-- update user set credits = credits - (select sum(n.cost_per_hour) from user_node un join node n on un.node_id = n.id where un.user_id = user.id);

create table minetest_server(
    id varchar(36) primary key not null, -- uuid
    user_node_id varchar(36) not null references user_node(id),
    name varchar(64) not null, -- name of the server, used for dns CNAME
    created bigint not null -- creation time in `time.Now().Unix()`
);