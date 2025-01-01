
export const get_user_settings = () => protected_fetch(`api/profile/settings`);

export const set_user_setting = (key, value) => protected_fetch(`api/profile/settings/${key}`, {
    method: "PUT",
    body: value
});

export const delete_user_setting = key => protected_fetch(`api/profile/settings/${key}`, {
    method: "DELETE"
});
