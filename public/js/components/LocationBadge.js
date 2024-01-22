import { country_map, flag_map } from "../util/country.js";

export default {
    props: ["location"],
    data: function() {
        return {
            flag_map,
            country_map
        };
    },
    template: /*html*/`
        <span class="" :title="country_map[location]">{{flag_map[location]}}</span>
        &nbsp;
    `
};