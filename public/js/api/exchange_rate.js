import { protected_fetch } from "./protected_fetch.js";

export const get_exchange_rates = () => protected_fetch(`api/exchange_rate`);
