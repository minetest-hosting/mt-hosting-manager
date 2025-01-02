
create table coupon(
    id varchar(36) primary key not null, -- uuid
    name varchar(128) not null,
    code varchar(32) not null unique,
    valid_from bigint not null, -- valid from time in `time.Now().Unix()`
    valid_until bigint not null, -- expiration time in `time.Now().Unix()`
    value bigint not null default 0 -- coupon value in euro-cent
);

create table redeemed_coupon(
    coupon_id varchar(36) not null references coupon(id) on delete cascade,
    user_id varchar(36) not null references public.user(id) on delete cascade,
    timestamp bigint not null, -- time in `time.Now().Unix()`
    primary key (coupon_id, user_id)
);