import { protected_fetch } from "./protected_fetch.js";

export const get_currency_info = () => protected_fetch(`api/currency`);
