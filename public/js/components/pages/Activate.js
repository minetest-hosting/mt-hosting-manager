import CardLayout from "../layouts/CardLayout.js";
import { activate } from "../../api/activation.js";
import { login } from "../../service/login.js";

export default {
    props: ["userid", "code"],
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "envelope", name: "Activate account", link: `/activate/${this.userid}/${this.code}`
            }],
            password1: "",
            password2: "",
            error: false
		};
	},
    methods: {
        activate: function() {
            activate({
                user_id: this.userid,
                activation_code: this.code,
                new_password: this.password1
            })
            .then(r => r.json())
            .then(resp => {
                if (resp.code && resp.code >= 400) {
                    this.error = true;
                } else {
                    return login({
                        mail: resp.mail,
                        password: this.password1
                    })
                    .then(() => this.$router.push("/profile"));
                }
            });
        }
    },
    computed: {
        password_valid: function() {
            if (this.password1 != this.password2) {
                return false;
            }
            if (this.password1.length < 8) {
                return false;
            }
            return true;
        }
    },
	template: /*html*/`
	<card-layout title="Activate account" icon="envelope" :breadcrumb="breadcrumb">
        <table class="table table-condensed table-striped">
            <tbody>
                <tr>
                    <td>User-ID</td>
                    <td>{{userid}}</td>
                </tr>
                <tr>
                    <td>Code</td>
                    <td>{{code}}</td>
                </tr>
                <tr>
                    <td>New password</td>
                    <td>
                        <input class="form-control" type="password" v-model="password1"/>
                    </td>
                </tr>
                <tr>
                    <td>New password (repeat)</td>
                    <td>
                        <input class="form-control" type="password" v-model="password2"/>
                        <div class="alert alert-info">
                            <i class="fa-solid fa-circle-info"></i>
                            <b>Note:</b> Minimum password length is 8 chars
                        </div>
                    </td>
                </tr>
                <tr>
                    <td>Action</td>
                    <td>
                        <button class="btn btn-primary w-100" :disabled="!password_valid || error" v-on:click="activate">
                            Activate
                        </button>
                        <div class="alert alert-danger" v-if="error">
                            <i class="fa-solid fa-triangle-exclamation"></i>
                            An error occured during activation
                        </div>
                    </td>
                </tr>
            </tbody>
        </table>
	</card-layout>
	`
};
