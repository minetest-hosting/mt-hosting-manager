
import store from '../store/exchange_rate.js';

export default {
    props: ["modelValue"],
    data: () => store,
    template: /*html*/`
    <select class="form-control" :value="modelValue" v-on:input="this.$emit('update:modelValue', $event.target.value)">
        <option v-for="rate in rates" :value="rate.currency">{{rate.display_name}} [{{rate.display_prefix}}]</option>
    </select>
    `
};
