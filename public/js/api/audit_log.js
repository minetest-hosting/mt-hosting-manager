import { protected_fetch } from "./protected_fetch.js";

export const search_audit_logs = s => protected_fetch(`api/audit_log`, {
    method: "POST",
    body: JSON.stringify(s)
});
