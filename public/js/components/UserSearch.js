import { search_user } from "../api/user.js";
import debounce from "../util/debounce.js";

export default {
    props: ["modelValue"],
    data: function() {
        return {
            show_modal: false,
            mail_like: this.modelValue ? this.modelValue.mail : "",
            busy: false,
            users: [],
            user: null
        };
    },
    created: function() {
        this.search();
    },
    methods: {
        select_user: function(user) {
            this.$emit("update:modelValue", user);
            this.user = user;
            this.show_modal = false;
        },
        search: debounce(function() {
            if (this.mail_like == "") {
                this.users = [];
                return;
            }
            this.busy = true;
            search_user({
                mail_like: `%${this.mail_like}%`,
                limit: 10
            })
            .then(l => {
                this.busy = false;
                this.users = l;
            });
        }, 250)
    },
    watch: {
        "mail_like": "search"
    },
    template: /*html*/`
    <div>
        <div class="input-group">
            <input class="form-control" readonly="true" :value="user && user.mail"/>
            <button class="btn btn-outline-secondary" v-on:click="select_user(null)">
                <i class="fa fa-trash"></i>
            </button>
            <button class="btn btn-outline-secondary" v-on:click="show_modal = true">
                <i class="fa fa-user"></i>
            </button>
        </div>
        <div class="modal show" style="display: block;" tabindex="-1" v-show="show_modal">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <h1 class="modal-title fs-5">
                            Search user
                            <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                        </h1>
                        <button type="button" class="btn-close" v-on:click="show_modal = null"></button>
                    </div>
                    <div class="modal-body">
                        <input type="text" class="form-control" v-model="mail_like"/>
                        <table class="table table-condensed table-striped">
                            <thead>
                                <tr>
                                    <th>Name</th>
                                    <th>Mail</th>
                                    <th>Select</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="user in users" :key="user.id" v-bind:class="{'table-success':modelValue && modelValue.id == user.id}">
                                    <td>{{user.name}}</td>
                                    <td>{{user.mail}}</td>
                                    <td>
                                        <i class="fa-regular fa-square-check" v-on:click="select_user(user)"></i>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    </div>
    `
};