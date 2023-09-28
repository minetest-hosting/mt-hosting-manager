import CardLayout from "../layouts/CardLayout.js";
import ServerLink from "../ServerLink.js";
import NodeLink from "../NodeLink.js";

export default {
	components: {
		"card-layout": CardLayout,
		"node-link": NodeLink,
		"server-link": ServerLink
	},
	data: function() {
		return {
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "question", name: "Help", link: "/help"
			}]
		};
	},
	template: /*html*/`
	<card-layout title="Help" icon="question" :breadcrumb="breadcrumb">
		<h4>Nodes and servers</h4>
		<p>
			This section explains the relationship between <i class="fa fa-server"></i> nodes and <i class="fa fa-list"></i> servers
		</p>
		<p>
			A node is a container for several (minetest) servers.
			From a financial standpoint only the node is billed, you <i>could</i> host as many minetest-servers on it as you like
			but the limiting factor will always be the available RAM and CPU resources.
			There are several types of nodes, each with their own resources.
		</p>
		<p>
			<b>Note:</b> ports (usually 30000+ for minetest) can only be used once on each node (see below).
		</p>
		<p>
			Below is an example how multiple servers are hosted on several nodes
		</p>
		<p>
			<ul>
				<li>
					<i class="fa fa-server"></i> Node 1, Side-projects, tests
					<ul>
						<li>
							<i class="fa fa-list"></i> My nodecore server (port 30000)
						</li>
						<li>
							<i class="fa fa-list"></i> My exile dev server (port 30001)
						</li>
					</ul>
				</li>
				<li>
					<i class="fa fa-server"></i> Node 2, Serious stuff
					<ul>
						<li>
							<i class="fa fa-list"></i> Nodecore skyhell main server (port 30000)
						</li>
						<li>
							<i class="fa fa-list"></i> Mineclone anarchy server (port 30001)
						</li>
					</ul>
				</li>
			</ul>
		</p>
		<p>
			The management of the servers can be facilitated with a web-interface
		</p>
		<hr>
		<h4>Payment, refunds and billing</h4>
		<p>
			Payments are deducted from your current balance, you can top that up at any time in the "Finance" page.
		</p>
		<p>
			<b>Note:</b> All local payments are done and displayed in the &euro; currency
		</p>
		<p>
			Refunds are possible for each payment (up to the available balance).
			Used up funds can't be refunded.
		</p>
		<p>
			Billing of nodes is done 24 hours in advance (see the "Next billing cycle" entry in the Node-details).
		</p>
	</card-layout>
	`
};
