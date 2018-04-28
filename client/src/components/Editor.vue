<template>
  <div>
    TODO
  </div>
</template>

<script>
import auth from '../auth'

export default {
  created: function () {
    var headers = auth.getAuthHeaders()
    this.$http.get('/api/activities', { headers: headers }).then(response => {
      console.log(response.body)
    }, response => {
      this.error = response.statusText
    })
  }
}
</script>
