import { protected_fetch } from "./protected_fetch.js";

export const get_all = () => protected_fetch(`api/job`);

export const retry = job => protected_fetch(`api/job/${job.id}`, {
    method: "POST"
});

export const remove = job => protected_fetch(`api/job/${job.id}`, {
    method: "DELETE"
});