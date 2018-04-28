<template>
  <div>
    App needs Strava permissions to see private activities and modify all activities.
  </div>
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
  }
}
</script>
