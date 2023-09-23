import { protected_fetch } from "./protected_fetch.js";

export const search_user = s => protected_fetch(`api/user/search`, {
    method: "POST",
    body: JSON.stringify(s)
});
