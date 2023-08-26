import { protected_fetch } from "./protected_fetch.js";

export const get_profile = () => protected_fetch(`api/profile`);

export const update_profile = u => protected_fetch(`api/profile`, {
    method: "POST",
    body: JSON.stringify(u)
});
