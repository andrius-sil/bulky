
const CLIENT_ID = '24292'
const REDIRECT_URL = `http://${process.env.HOST}:${process.env.PORT}/`

const ACCESS_TOKEN = 'access_token'

export default {
  getAccessToken () {
    return localStorage.getItem(ACCESS_TOKEN)
  },

  isAuthenticated () {
    return !!this.getAccessToken()
  },

  stravaRequestAccess () {
    var url = `https://www.strava.com/oauth/authorize?client_id=${CLIENT_ID}&response_type=code&redirect_uri=${REDIRECT_URL}&scope=view_private,write`
    window.location.href = url
  },

  stravaTokenExchange (context, authCode, redirectPage) {
    var payload = { code: authCode }
    context.$http.post('/api/login', payload).then(response => {
      localStorage.setItem(ACCESS_TOKEN, response.body.access_token)
      context.$router.replace(redirectPage)
    }, response => {
      context.error = response.statusText
    })
  },

  deauthenticate (context) {
    var payload = { access_token: localStorage.getItem(ACCESS_TOKEN) }
    context.$http.post('https://www.strava.com/oauth/deauthorize', payload).then(response => {
      context.$router.replace('/')
    }, response => {
      context.error = response.statusText
    })
    localStorage.removeItem(ACCESS_TOKEN)
  },

  getAuthHeaders () {
    return { 'Authorization': 'Bearer ' + this.getAccessToken() }
  }
}
