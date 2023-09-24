import { protected_fetch } from "./protected_fetch.js";

export const send_mail = (m, user_id) => protected_fetch(`api/mail/send/${user_id}`, {
    method: "POST",
    body: JSON.stringify(m)
});