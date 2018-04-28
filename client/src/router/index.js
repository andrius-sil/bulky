import Vue from 'vue'
import Router from 'vue-router'
import Editor from '@/components/Editor'
import Home from '@/components/Home'
import NotFound from '@/components/NotFound'

import auth from '../auth'

Vue.use(Router)

function requireNoAuth (to, from, next) {
  if (auth.isAuthenticated()) {
    next('/editor')
  } else {
    next()
  }
}

function requireAuth (to, from, next) {
  if (!auth.isAuthenticated()) {
    next('/')
  } else {
    next()
  }
}

// TODO: error page for non-existing routes
export default new Router({
  mode: 'history', // TODO: why query params won't work otherwise?
  routes: [
    {
      path: '/',
      component: Home,
      beforeEnter: requireNoAuth
    },
    {
      path: '/editor',
      component: Editor,
      beforeEnter: requireAuth
    },
    {
      path: '*',
      component: NotFound
    }
  ]
})
