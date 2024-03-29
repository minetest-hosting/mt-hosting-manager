-- registered user
create table public.user(
    id varchar(36) primary key not null, -- uuid
    state varchar(32) not null default 'ACTIVE',
    name varchar(128) not null, -- username
    hash varchar(256) not null default '', -- bcrypt hash for local users
    created bigint not null, -- creation time in `time.Now().Unix()`
    balance bigint not null default 0, -- current balance in euro-cent
    external_id varchar(64) not null, -- id on the external oauth provider
    currency varchar(16) not null default 'EUR', -- user preferred currency
    type varchar(32) not null, -- GITHUB, DISCORD
    role varchar(32) not null -- ADMIN / USER
);

create unique index user_name on public.user(name);

-- available node types to select
create table node_type(
    id varchar(36) primary key not null, -- uuid
    state varchar(32) not null default 'INACTIVE',
    order_id int not null default 1, -- order id
    provider varchar(32) not null, -- HETZNER
    server_type varchar(32) not null, -- provider server-type: cx11
    locations varchar(32) not null, -- location list
    name varchar(128) not null default '', -- name of the node
    description varchar(1024) not null default '', -- description of the node
    cpu_count int not null default 1, -- number of cpu's
    ram_gb int not null default 2, -- ram in gb
    disk_gb int not null default 10, -- disk size in gb
    dedicated boolean not null default false, -- dedicated machine
    daily_cost bigint not null default 0, -- daily cost in eurocents
    max_recommended_instances int not null default 2, -- max number of recommended minetest instances on this host
    max_instances int not null default 4 -- max number of allowed minetest instances on this host
);

-- a node set up by a user
create table user_node(
    id varchar(36) primary key not null, -- uuid
    user_id varchar(36) not null references public.user(id) on delete cascade,
    node_type_id varchar(36) not null references node_type(id) on delete cascade,
    location varchar(32) not null default 'nbg1',
    external_id varchar default '',
    created bigint not null, -- creation time in `time.Now().Unix()`
    valid_until bigint not null, -- validity ("payed" until) in `time.Now().Unix()`
    state varchar(32) not null default 'CREATED',
    name varchar(64) not null, -- name of the host, used for dns registration (A, AAAA record)
    alias varchar(256) not null, -- internal name, user-specified
    ipv4 varchar(32) not null,
    ipv6 varchar(128) not null,
    external_ipv4_dns_id varchar default '', -- hetzner dns record id
    external_ipv6_dns_id varchar default '', -- hetzner dns record id
    fingerprint varchar(512) not null default '' -- ssh fingerprint
);

create table minetest_server(
    id varchar(36) primary key not null, -- uuid
    user_node_id varchar(36) not null references user_node(id) on delete cascade,
    name varchar(64) not null, -- display name of the server
    dns_name varchar(256) not null, -- DNS name prefix
    admin varchar(32) not null default 'admin', -- admin name
    external_cname_dns_id varchar default '', -- hetzner dns record id
    custom_dns_name varchar default '', -- custom dns name (CNAME pointing to server-CNAME) optional
    port int not null default 30000, -- minetest server port
    ui_version varchar(16) not null default 'sha256:..', -- ui version to deploy
    jwt_key varchar(16) not null, -- jwt key
    created bigint not null, -- creation time in `time.Now().Unix()`
    state varchar(32) not null default 'CREATED'
);

create table backup_space(
    id varchar(36) primary key not null, -- uuid
    name varchar(64) not null, -- display name
    user_id varchar(36) not null references public.user(id) on delete cascade, -- belongs to user
    retention_days int not null default 30, -- backup retention in days
    created bigint not null, -- creation time in `time.Now().Unix()`
    valid_until bigint not null -- validity ("payed" until) in `time.Now().Unix()`
);

create table backup(
    id varchar(36) primary key not null, -- uuid
    state varchar(32) not null default 'CREATED', -- CREATED, PROGRESS, COMPLETE, ERROR
    passphrase varchar(128) not null, -- passphrase to encrypt tar-file with (`openssl enc -aes-256-cbc -pbkdf2 -pass pass:${passphrase}`)
    backup_space_id varchar(36) not null references backup_space(id) on delete cascade, -- belongs to backup_space
    minetest_server_id varchar(36) references minetest_server(id) on delete cascade,
    created bigint not null, -- creation time in `time.Now().Unix()`
    size bigint not null -- size of the backup in bytes
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
    data bytea -- json job data
);

create index job_type_state_started on job(type, state, started);

create table payment_transaction(
    id varchar(36) primary key not null, -- uuid
    type varchar(32) not null, -- WALLEE, COINBASE
    transaction_id varchar not null, -- external tx id / coinbase code
    payment_url varchar not null default '', -- url to payment service
    created bigint not null, -- creation time in `time.Now().Unix()`
    expires bigint not null, -- expiry time in `time.Now().Unix()`
    user_id varchar(36) not null references public.user(id) on delete cascade,
    amount bigint not null default 0, -- currency amount in eurocents
    amount_refunded bigint not null default 0, -- amount refunded from this transaction in eurocents
    state varchar(32) not null default 'PENDING' -- state of the transaction, PENDING, SUCCESS, ERROR
);

create index payment_transaction_user_id on payment_transaction(user_id);

create table audit_log(
    id varchar(36) primary key not null, -- uuid
    type varchar(64) not null, -- type of audit log
    timestamp bigint not null, -- time in `time.Now().Unix()`
    user_id varchar(36) not null references public.user(id) on delete cascade, -- user
    ip_address varchar(128), -- ip address (optional)
    user_node_id varchar(36), -- node (optional)
    minetest_server_id varchar(36), -- server (optional)
    payment_transaction_id varchar(36), -- payment (optional)
    amount bigint -- currency amount in euro-cent
);

create index audit_log_search on audit_log(type, timestamp, user_id);

create table exchange_rate(
    currency varchar(16) primary key not null, -- EUR, ETH, BTC
    rate varchar(64) not null, -- EUR * rate = $currency
    display_name varchar not null, -- "United States Dollar"
    digits int not null default 2, -- 0.01
    updated bigint not null -- update time in `time.Now().Unix()`
);
