
alter table job add backup_id varchar(36) references backup(id) on delete cascade;
alter table audit_log add backup_id varchar(36) references backup(id) on delete cascade;