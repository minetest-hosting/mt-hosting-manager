import Start from './components/pages/Start.js';
import Profile from './components/pages/Profile.js';
import NodeTypes from './components/pages/NodeTypes.js';
import MTServers from './components/pages/MTServers.js';
import UserNodes from './components/pages/UserNodes.js';
import Jobs from './components/pages/Jobs.js';

export default [{
	path: "/", component: Start
},{
	path: "/profile", component: Profile
},{
	path: "/nodes", component: UserNodes
},{
	path: "/mtservers", component: MTServers
},{
	path: "/jobs", component: Jobs,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/nodetypes", component: NodeTypes,
	meta: { requiredRole: "ADMIN" }
}];
