import Vue from 'vue'
import App from './App'
import VueResource from 'vue-resource'
import router from './router'

import 'bulma/css/bulma.css'
import 'font-awesome/css/font-awesome.css'

Vue.use(VueResource)

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
