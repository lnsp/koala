import Vue from 'vue';
import App from './App.vue';
import router from './router';
import axios from 'axios';
import { createOidcAuth, SignInType, LogLevel } from 'vue-oidc-client';

function setupVue() {
  new Vue({
    router,
    render: h => h(App)
  }).$mount('#app');
}

axios.get('/', { baseURL: process.env.VUE_APP_API_BASEURL }).then((value) => {
  Vue.prototype.$baseURL = process.env.VUE_APP_API_BASEURL;
  if (value.data.security === 'oidc') {
    let oidc = createOidcAuth('main', SignInType.Popup, location.protocol + '//' + location.host + '/#/', {
      authority: value.data.provider,
      client_id: value.data.clientId,
      response_type: 'id_token token',
      scope: 'openid profile email',
      popupWindowTarget: '_blank',
      popupWindowFeatures: '',
    }, console, LogLevel.Debug);
    Vue.prototype.$oidc = oidc;
    oidc.useRouter(router);
    oidc.startup().then(() => setupVue());
  } else {
    setupVue();
  }
})
