<template>
  <div id="app" class="container">
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <h1 class="title">Bulky</h1>
      </div>

      <div class="navbar-menu">
        <div class="navbar-end">
          <button class="button is-link is-outlined is-small" v-show="!isAuthenticated() && !isHomepage" @click="handleLogin()">Authenticate with Strava</button>
          <button class="button is-link is-outlined is-small" v-show="isAuthenticated()" @click="handleLogout()">Revoke Strava Access</button>
        </div>
      </div>
    </nav>

    <router-view/>

    <footer class="footer">
      <div class="container">
        <div class="content has-text-centered">
          <p>
            Bulk Editor for Strava Activities by <i>Andrius Silinskas</i>.
          </p>
        </div>
      </div>
    </footer>
    </div>
</template>

<script>
import auth from './auth'

export default {
  name: 'App',

  methods: {
    isAuthenticated () {
      return auth.isAuthenticated()
    },
    handleLogin () {
      auth.stravaRequestAccess()
    },
    handleLogout () {
      auth.deauthenticate(this)
    },
    isHomepage () {
      return this.$route.path === '/'
    }
  }
}
</script>

<style>
#app {
  margin-top: 10px;
}
</style>
