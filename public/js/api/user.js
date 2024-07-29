import { protected_fetch } from "./protected_fetch.js";

export const search_user = s => protected_fetch(`api/user/search`, {
    method: "POST",
    body: JSON.stringify(s)
});

export const save_user = u => protected_fetch(`api/user/${u.id}`, {
    method: "POST",
    body: JSON.stringify(u)
});

export const get_users = () => protected_fetch(`api/user`);

export const get_user_by_id = id => protected_fetch(`api/user/${id}`);
