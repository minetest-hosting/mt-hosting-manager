-- registered user
create table user(
    id varchar(36) primary key not null, -- uuid
    state varchar(32) not null default 'ACTIVE',
    name varchar(128) not null, -- username
    mail varchar(128) not null, -- email
    created bigint not null, -- creation time in `time.Now().Unix()`
    external_id varchar(64) not null, -- id on the external oauth provider
    type varchar(32) not null, -- GITHUB, DISCORD
    role varchar(32) not null -- ADMIN / USER
);

create unique index user_mail on user(mail);

-- available node types to select
create table node_type(
    id varchar(36) primary key not null, -- uuid
    state varchar(32) not null default 'INACTIVE',
    order_id int not null default 1, -- order id
    provider varchar(32) not null, -- HETZNER
    server_type varchar(32) not null, -- provider server-type: cx11
    name varchar(128) not null default '', -- name of the node
    description varchar(1024) not null default '', -- description of the node
    monthly_cost varchar not null default '', -- "5.20"
    max_months int not null default '', -- 12
    max_recommended_instances int not null default 2, -- max number of recommended minetest instances on this host
    max_instances int not null default 4 -- max number of allowed minetest instances on this host
);

-- default node types
INSERT INTO node_type VALUES('0b71901c-9fe7-4a49-9431-e8ce5981310c','ACTIVE',0,'HETZNER','cx11','SMALL1','Small, versatile node, suited for 1 or 2 minetest servers','8.5',24,2,4);
INSERT INTO node_type VALUES('37d9f80b-8a4e-4c22-bd7a-65ad23ae1fa4','ACTIVE',5,'HETZNER','cx21','MEDIUM1','Medium node for average servers and mod-sets','12.5',12,3,6);
INSERT INTO node_type VALUES('fedbbf78-ef43-4fa6-9f1c-b24180c93ac3','ACTIVE',10,'HETZNER','cx41','LARGE1','Larger node for heavier workloads','17.5',12,5,10);

-- a node set up by a user
create table user_node(
    id varchar(36) primary key not null, -- uuid
    user_id varchar(36) not null references user(id) on delete restrict,
    node_type_id varchar(36) not null references node_type(id) on delete restrict,
    external_id varchar default '',
    created bigint not null, -- creation time in `time.Now().Unix()`
    expires bigint not null, -- expiration time in `time.Now().Unix()`
    state varchar(32) not null default 'CREATED',
    name varchar(64) not null, -- name of the host, used for dns registration (A, AAAA record)
    alias varchar(256) not null, -- internal name, user-specified
    ipv4 varchar(32) not null,
    ipv6 varchar(128) not null,
    load_percent int not null default 0,
    disk_size int not null default 0,
    disk_used int not null default 0,
    memory_size int not null default 0,
    memory_used int not null default 0
);

-- update user set credits = credits - (select sum(n.cost_per_hour) from user_node un join node n on un.node_id = n.id where un.user_id = user.id);

create table minetest_server(
    id varchar(36) primary key not null, -- uuid
    user_node_id varchar(36) not null references user_node(id) on delete restrict,
    name varchar(64) not null, -- name of the server, used for dns CNAME
    created bigint not null, -- creation time in `time.Now().Unix()`
    state varchar(32) not null default 'CREATED'
);

create table job(
    id varchar(36) primary key not null, -- uuid
    type varchar(32) not null, -- NODE_SETUP, NODE_DESTROY, SERVER_SETUP, SERVER_DESTROY
    state varchar(32) not null default 'CREATED', -- CREATED, RUNNING, DONE_SUCCESS, DONE_FAILURE
    started bigint not null, -- start time in `time.Now().Unix()`
    finished bigint not null, -- end time in `time.Now().Unix()`
    user_node_id varchar(36) references user_node(id) on delete cascade,
    minetest_server_id varchar(36) references minetest_server(id) on delete cascade,
    progress_percent float not null default 0,
    message varchar not null default '',
    data blob -- json job data
);

create index job_type_state_started on job(type, state, started);

create table payment_transaction(
    id varchar(36) primary key not null, -- uuid
    transaction_id varchar not null, -- external tx id
    created bigint not null, -- creation time in `time.Now().Unix()`
    node_type_id varchar(36) not null references node_type(id) on delete restrict,
    months int not null
);