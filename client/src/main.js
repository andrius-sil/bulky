import Vue from 'vue'
import App from './App'
import VueResource from 'vue-resource'
import router from './router'

import 'bulma/css/bulma.css'
import 'font-awesome/css/font-awesome.css'

import * as VueGoogleMaps from 'vue2-google-maps'

Vue.use(VueResource)

Vue.use(VueGoogleMaps, {
  load: {
    key: process.env.GOOGLE_MAPS_API_KEY,
    libraries: 'places'
  }
})

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
