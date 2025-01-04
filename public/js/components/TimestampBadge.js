import format_duration from "../util/format_duration.js";
import format_time from "../util/format_time.js";

export default {
    props: ["timestamp", "show_duration"],
    methods: {
        format_duration,
        format_time
    },
    computed: {
        duration: function() {
            return (Date.now()/1000) - this.timestamp;
        },
        str: function() {
            let s = this.format_time(this.timestamp);
            if (this.show_duration) {
                s += " (";
                if (this.duration < 0) {
                    s += "in ";
                }
                s += this.format_duration(this.duration);
                if (this.duration >= 0) {
                    s += " ago";
                }
                s += ")";
            }
            return s;
        }
    },
    template: /*html*/ `{{str}}`
};
