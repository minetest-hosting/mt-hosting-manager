
export const send_activation = sar => fetch(`api/send_activation`, {
    method: "POST",
    body: JSON.stringify(sar)
});

export const activate = ar => fetch(`api/activate`, {
    method: "POST",
    body: JSON.stringify(ar)
});