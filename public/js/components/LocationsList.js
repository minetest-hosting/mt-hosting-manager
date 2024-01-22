import LocationBadge from "./LocationBadge.js";

export default {
    props: ["locations"],
    components: {
        "location-badge": LocationBadge
    },
    computed: {
        location_list: function() {
            return this.locations.split(",");
        }
    },
    template: /*html*/`
        <location-badge :location="location" :key="location" v-for="location in location_list"/>
    `
};