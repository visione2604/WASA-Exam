import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginPage/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import SearchView from '../views/SearchView.vue'
import ConvView from '../views/ConvView.vue'
import GroupCreateView from '../views/GroupCreateView.vue'
import GroupEditView from '../views/GroupEditView.vue'
import UsersView from '../views/UsersView.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', redirect: '/login'},
		{path: '/login', component: LoginView},
		{path: '/home', component: HomeView, meta: { requiresAuth: true }},
		{path: '/profile', component: ProfileView, meta: { requiresAuth: true }},
		{path: '/users', component: UsersView, meta: { requiresAuth: true }},
		{path: '/search', component: SearchView, meta: { requiresAuth: true }},
		{path: '/conversations/:conversationId', component: ConvView, meta: { requiresAuth: true }},
		{path: '/groups/create', component: GroupCreateView, meta: { requiresAuth: true }},
		{path: '/groups/:groupId/edit', component: GroupEditView, meta: { requiresAuth: true }},
	]
})

router.beforeEach((to, from, next) => {
	const isAuthenticated = localStorage.getItem('token')
	
	if (to.meta.requiresAuth && !isAuthenticated) {
		next('/login')
	} else if (to.path === '/login' && isAuthenticated) {
		next('/home')
	} else {
		next()
	}
})

export default router
