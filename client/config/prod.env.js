'use strict'
module.exports = {
  NODE_ENV: '"production"',
  GOOGLE_MAPS_API_KEY: JSON.stringify(process.env.GOOGLE_MAPS_API_KEY),
  HOST: JSON.stringify(process.env.PRODUCTION_HOST),
  PORT: 8080
}
