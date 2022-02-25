import { createApp } from 'vue'
import { createStore } from 'vuex'
import App from './App.vue'
import VueAxios from 'vue-axios'
import axios from 'axios'

// setup store
const store = createStore({
  state() {
    return {
      records: [],
      zones: [],
    }
  },
  mutations: {
    addRecord (state, { zone, record }) {
      record.zone = zone
      state.records.push(record)
    },
    dropRecord (state, { zone, record }) {
      record.zone = zone
      state.records = state.records.filter(r => r !== record)
    },
    setRecords (state, { zone, records }) {
      state.records = state.records.filter(r => r.zone !== zone)
      records.forEach(r => {
        r.zone = zone
        state.records.push(r)
      })
    },
    setZones (state, { zones }) {
      state.zones = zones
    }
  }
})

// configure app
const app = createApp(App)
app.use(VueAxios, axios.create({
  baseURL: process.env.VUE_APP_API_BASEURL,
}))
app.use(store)
app.provide('axios', app.config.globalProperties.axios)
app.mount('#app')