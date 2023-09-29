
export const set_password = spr => fetch(`api/set_password`, {
    method: "POST",
    body: JSON.stringify(spr)
});