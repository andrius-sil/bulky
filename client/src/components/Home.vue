<template>
  <section class="hero is-primary is-medium">
    <div class="hero-body">
      <div class="container has-text-centered">
        <h1 class="subtitle">
          <strong>App needs Strava permissions to see private activities and modify all activities.</strong>
        </h1>
        <a @click="handleLogin()"><img src="../assets/btn_strava_connectwith_orange.png"></a>
      </div>
    </div>
  </section>
</template>

<script>
import auth from '../auth'

export default {
  created: function () {
    // Strava auth code after redirect.
    var authCode = this.$route.query.code
    if (authCode != null) {
      auth.stravaTokenExchange(this, authCode, '/editor')
    }
  },

  methods: {
    handleLogin () {
      auth.stravaRequestAccess()
    }
  }
}
</script>
