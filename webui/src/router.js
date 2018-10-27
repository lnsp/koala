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
                'rootAPI': process.env.VUE_APP_ROOT_API,
            }
        },
        {
            path: '/auth',
            name: 'authenticate',
            component: Authentication,
            props: {
                'rootAPI': process.env.VUE_APP_ROOT_API,
            }
        }
    ],
});