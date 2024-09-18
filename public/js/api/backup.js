import { protected_fetch } from "./protected_fetch.js";

export const get_by_backup_space_id = id => protected_fetch(`api/backup_space/${id}/backup`);

export const get_by_id = id => protected_fetch(`api/backup/${id}`);

export const create = b => protected_fetch(`api/backup`, {
    method: "POST",
    body: JSON.stringify(b)
});

export const remove = b => protected_fetch(`api/backup/${b.id}`, {
    method: "DELETE"
});