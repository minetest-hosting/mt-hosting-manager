import { protected_fetch } from "./protected_fetch.js";

export const get_all = () => protected_fetch(`api/backup_space`);
export const get_by_id = id => protected_fetch(`api/backup_space/${id}`);

export const create = bs => protected_fetch(`api/backup_space`, {
    method: "POST",
    body: JSON.stringify(bs)
});

export const update = bs => protected_fetch(`api/backup_space/${bs.id}`, {
    method: "POST",
    body: JSON.stringify(bs)
});

export const remove = bs => protected_fetch(`api/backup_space/${bs.id}`, {
    method: "DELETE"
});