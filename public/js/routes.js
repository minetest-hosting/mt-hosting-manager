import Start from './components/pages/Start.js';
import Profile from './components/pages/Profile.js';
import Login from './components/pages/Login.js';
import NodeTypes from './components/pages/NodeTypes.js';
import NodeTypeDetail from './components/pages/NodeTypeDetail.js';
import MTServers from './components/pages/MTServers.js';
import MTServerDetail from './components/pages/MTServerDetail.js';
import UserNodes from './components/pages/UserNodes.js';
import UserNodesDetail from './components/pages/UserNodesDetail.js';
import Jobs from './components/pages/Jobs.js';
import Finance from './components/pages/Finance.js';
import FinanceDetail from './components/pages/FinanceDetail.js';

export default [{
	path: "/", component: Start
},{
	path: "/login", component: Login
},{
	path: "/profile", component: Profile,
	meta: { loggedIn: true }
},{
	path: "/finance", component: Finance,
	meta: { loggedIn: true }
},{
	path: "/finance/detail/:id", component: FinanceDetail,
	meta: { loggedIn: true }
},{
	path: "/nodes", component: UserNodes,
	meta: { loggedIn: true }
},{
	path: "/nodes/:id", component: UserNodesDetail,
	meta: { loggedIn: true }
},{
	path: "/mtservers", component: MTServers,
	meta: { loggedIn: true }
},{
	path: "/mtservers/:id", component: MTServerDetail,
	meta: { loggedIn: true }
},{
	path: "/jobs", component: Jobs,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/node_types", component: NodeTypes,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/node_types/:id", component: NodeTypeDetail,
	meta: { requiredRole: "ADMIN" }
}];
