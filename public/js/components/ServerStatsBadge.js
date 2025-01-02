import { get_stats } from "../api/mtserver.js";

export default {
    props: ["id"],
    data: function() {
        return {
            handle: null,
            stats: null
        };
    },
    mounted: function() {
        this.update();
        this.handle = setInterval(() => this.update(), 2000);
    },
	beforeUnmount: function() {
		clearInterval(this.handle);
	},
    methods: {
        update: async function() {
            this.stats = await get_stats(this.id);
        }
    },
    computed: {
        signalColor: function() {
            if (this.stats.max_lag < 200) return "green";
            if (this.stats.max_lag < 500) return "yellow";
            return "red";
        },
        hour: function() {
            return Math.floor(this.stats.time_of_day * 24);
        },
        minute: function() {
            const min = Math.floor(((this.stats.time_of_day * 24) - this.hour) * 60);
            return min >= 10 ? min : "0" + min;
        }
    },
    template: /*html*/`
    <span v-if="!stats">
        <i class="fa fa-spin fa-spinner"></i>
    </span>
    <span class="badge bg-warning" v-else-if="stats.maintenance">
        Maintenance
    </span>
    <span class="badge bg-warning" v-else-if="stats.uptime == 0">
        Stopped
    </span>
    <span v-else>
    <i class="fa-solid fa-signal" v-bind:style="{'color': signalColor}"></i>
        {{ Math.floor(stats.max_lag*1000) }} ms
        <i class="fa-solid fa-users"></i>
        {{ stats.player_count }}
        <i class="fa-solid fa-clock"></i>
        {{hour}}:{{minute}}
        <i class="fa-solid fa-sun" style="color: yellow;" v-if="hour >= 6 && hour < 18"></i>
        <i class="fa-solid fa-moon" style="color: lightblue;" v-else></i>
    </span>
    `
};