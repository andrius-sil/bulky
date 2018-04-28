<template>
  <div id="app">
    <div>
      <button v-show="!isAuthenticated()" @click="handleLogin()">Authenticate with Strava</button>
      <button v-show="isAuthenticated()" @click="handleLogout()">Log out</button>
    </div>
    <router-view/>
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
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
