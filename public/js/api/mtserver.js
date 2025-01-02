import { protected_fetch } from "./protected_fetch.js";

export const get_all = () => protected_fetch(`api/mtserver`);

export const get_by_id = id => protected_fetch(`api/mtserver/${id}`);

export const get_latest_job = id => protected_fetch(`api/mtserver/${id}/job`);

export const get_stats = id => protected_fetch(`api/mtserver/${id}/stats`);

export const create = (s, restore_from) => protected_fetch(`api/mtserver${restore_from ? '?restore_from=' + restore_from : ''}`, {
    method: "POST",
    body: JSON.stringify(s)
});

export const create_validate = s => protected_fetch(`api/mtserver/validate`, {
    method: "POST",
    body: JSON.stringify(s)
});


export const update = s => protected_fetch(`api/mtserver/${s.id}`, {
    method: "POST",
    body: JSON.stringify(s)
});

export const setup = s => protected_fetch(`api/mtserver/${s.id}/setup`, {
    method: "POST",
});


export const remove = s => protected_fetch(`api/mtserver/${s.id}`, {
    method: "DELETE"
});