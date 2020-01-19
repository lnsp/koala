import Vue from 'vue';
import Router from 'vue-router';
import ControlPanel from './views/ControlPanel.vue';
import Authentication from './views/Authentication.vue';

Vue.use(Router);

export default new Router({
    routes: [
        {
            path: '/',
            name: 'controlPanel',
            component: ControlPanel,
            props: {
                'apiBaseURL': process.env.VUE_APP_API_BASEURL,
            }
        },
        {
            path: '/auth',
            name: 'authenticate',
            component: Authentication,
            props: {
                'apiBaseURL': process.env.VUE_APP_API_BASEURL,
            }
        }
    ],
});