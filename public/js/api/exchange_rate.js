import { protected_fetch } from './util.js';

export const get_exchange_rates = () => protected_fetch(`api/exchange_rate`);