
alter table backup add column user_id varchar(36);
update backup set user_id = (select user_id from backup_space where id = backup.backup_space_id);
alter table backup alter column user_id set not null;
alter table backup
    add constraint backup_user_id_fkey
    foreign key (user_id)
    references public.user(id)
    on delete cascade;

alter table backup drop column backup_space_id;
drop table backup_space;

alter table backup add column expires bigint;
update backup set expires = created + (3600 * 24 * 365);
alter table backup alter column expires set not null;
