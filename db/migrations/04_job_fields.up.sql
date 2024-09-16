alter table job rename column started to created;
alter table job add column next_run bigint not null default 0;
alter table job add column step int not null default 0;
