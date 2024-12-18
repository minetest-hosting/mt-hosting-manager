import { protected_fetch } from "./protected_fetch.js";

export const get_all = () => protected_fetch(`api/node`);
export const get_by_id = id => protected_fetch(`api/node/${id}`);
export const get_mtservers_by_nodeid = id => protected_fetch(`api/node/${id}/mtservers`);
export const get_latest_job = id => protected_fetch(`api/node/${id}/job`);
export const get_stats = id => protected_fetch(`api/node/${id}/stats`);

export const search = s => protected_fetch(`api/node/search`, {
    method: "POST",
    body: JSON.stringify(s)
});

export const create = n => protected_fetch(`api/node`, {
    method: "POST",
    body: JSON.stringify(n)
});

export const update = n => protected_fetch(`api/node/${n.id}`, {
    method: "POST",
    body: JSON.stringify(n)
});

export const remove = n => protected_fetch(`api/node/${n.id}`, {
    method: "DELETE"
});