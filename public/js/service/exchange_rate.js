import store from "../store/exchange_rate.js";
import { get_exchange_rates } from "../api/exchange_rate.js";

export const fetch_exchange_rates = () => get_exchange_rates().then(r => store.rates = r);