import { protected_fetch } from "./utiljs";

export const get_all = () => protected_fetch(`api/mtserver`);

export const create = s => protected_fetch(`api/mtserver`, {
    metod: "POST",
    body: JSON.stringify(s)
});

export const update = s => protected_fetch(`api/mtserver/${s.id}`, {
    metod: "POST",
    body: JSON.stringify(s)
});

export const remove = s => protected_fetch(`api/mtserver/${s.id}`, {
    metod: "DELETE"
});