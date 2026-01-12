import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'

import './assets/dashboard.css'
import './assets/main.css'

// Ensure fresh session on each load: clear auth persistence
try {
	localStorage.removeItem('token');
	localStorage.removeItem('username');
	localStorage.removeItem('userId');
	localStorage.removeItem('userPhoto');
} catch (e) {}
try {
	if (axios && axios.defaults && axios.defaults.headers && axios.defaults.headers.common) {
		delete axios.defaults.headers.common['Authorization'];
	}
} catch (e) {}

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.use(router)
app.mount('#app')
