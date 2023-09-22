
export default {
    props: ["eurocents"],
    template: /*html*/`
    <span>
        &euro; {{eurocents/100}}
    </span>
    `
};