import { protected_fetch } from "./protected_fetch.js";

export const get_overview = user_id => protected_fetch(`api/overview/${user_id}`);
