<template>
  <div>
    <div class="level">
      <div class="level-left">
        <div class="level-item">
          <p>from</p>
        </div>
        <div class="level-item">
          <datepicker wrapper-class="control" input-class="input" v-model="startDate"></datepicker>
        </div>

        <div class="level-item"></div>

        <div class="level-item">
          <p>to</p>
        </div>
        <div class="level-item">
          <datepicker wrapper-class="control" input-class="input" v-model="endDate"></datepicker>
        </div>
      </div>

      <div class="level-right">
        <div class="level-item">
          <button class="button is-primary" :disabled="isUpdateDisabled()" @click="updateActivities()">Update Activities</button>
        </div>
      </div>
    </div>

    <div class="level">
      <div class="level-left">
        <div class="level-item">
          <p>start</p>
        </div>
        <div class="level-item">
          <GmapAutocomplete class="input" placeholder="Enter home address" @place_changed="setStartPlace"></GmapAutocomplete>
        </div>

        <div class="level-item"></div>

        <div class="level-item">
          <p>end</p>
        </div>
        <div class="level-item">
          <GmapAutocomplete class="input" placeholder="Enter work address" @place_changed="setEndPlace"></GmapAutocomplete>
        </div>

        <div class="level-item"></div>

        <div class="level-item">
          <p>proximity offset</p>
        </div>
        <div class="level-item">
          <input v-model="proximityMeters" class="slider" type="range" step="50" min="0" max="1000">
        </div>
        <div class="level-item">
          <p>{{ proximityMeters }} metres</p>
        </div>

        <div class="level-item"></div>

        <div class="level-item">
          <label class="checkbox">
            <input type="checkbox" v-model="autoSelectCommutes">
            Auto select commutes
          </label>
        </div>
      </div>
    </div>

    <div class="container">
      <table class="table is-fullwidth">
        <tr>
          <th>Date</th>
          <th>Name</th>
          <th>Distance (km)</th>
          <th>
            <input type="checkbox" v-model="selectAllCommutes">
            Commute
          </th>
          <th>
            <input type="checkbox" v-model="selectAllPrivate">
            Private
          </th>
        </tr>
        <tr v-for="activity in activities" :key="activity.Id" :class="getClassForActivity(activity.Id)">
          <td>{{ formatDate(activity.Start_date) }}</td>
          <td><a :href="activity.Url">{{ activity.Name }}</a></td>
          <td>{{ formatDistance(activity.Distance) }}</td>
          <td>
            <input type="checkbox" v-model="selectedCommute" :value="activity.Id" number>
          </td>
          <td>
            <input type="checkbox" v-model="selectedPrivate" :value="activity.Id" number>
          </td>
        </tr>
      </table>

      <div class="has-text-centered" v-if="activities.length === 0">
        <p class="no-activities">
          <strong>No activities for selected dates.</strong>
        </p>
      </div>

      <div class="has-text-centered" v-if="activities.length === 30">
        <p>
          <strong>Only first 30 activities are displayed. Try a smaller date range to edit the rest of activities.</strong>
        </p>
      </div>
    </div>

  </div>
</template>

<script>
import haversineDistance from 'haversine-distance'
import dateFormat from 'dateformat'
import Datepicker from 'vuejs-datepicker'

import 'bulma-slider/dist/bulma-slider.min.css'

import auth from '../auth'

// Today's date with time as 23:59:59.
function dateNow () {
  var now = new Date()
  now.setHours(23, 59, 59, 0)
  return now
}

// Date a week ago with time as 00:00:00.
function dateWeekAgo () {
  var now = new Date()
  now.setHours(0, 0, 0, 0)
  now.setDate(now.getDate() - 7)
  return now
}

function dateToEpochTimestamp (d) {
  return (d.getTime() / 1000).toFixed(0)
}

function arrayDiff (arr1, arr2) {
  var diff = arr1
    .filter(x => !arr2.includes(x))
    .concat(arr2.filter(x => !arr1.includes(x)))

  return diff
}

export default {
  data: function () {
    return {
      activities: [],

      selectedPrivate: [],
      fetchedPrivate: [],
      modifiedPrivate: [],

      selectedCommute: [],
      fetchedCommute: [],
      modifiedCommute: [],

      startDate: dateWeekAgo(),
      endDate: dateNow(),

      startPlaceGeometry: null,
      endPlaceGeometry: null,
      proximityMeters: 300,
      autoSelectCommutes: true
    }
  },

  components: {
    Datepicker
  },

  watch: {
    selectedPrivate: function (val) {
      this.modifiedPrivate = arrayDiff(this.selectedPrivate, this.fetchedPrivate)
    },
    selectedCommute: function (val) {
      this.modifiedCommute = arrayDiff(this.selectedCommute, this.fetchedCommute)
    },
    startDate: function (val) {
      this.fetchActivities()
    },
    endDate: function (val) {
      this.fetchActivities()
    },
    proximityMeters: function (val) {
      this.selectCommutes()
    },
    autoSelectCommutes: function (val) {
      this.selectCommutes()
    }
  },

  created: function () {
    this.fetchActivities()
  },

  methods: {
    fetchActivities: function () {
      var headers = auth.getAuthHeaders()
      var params = {
        after: dateToEpochTimestamp(this.startDate),
        before: dateToEpochTimestamp(this.endDate)
      }

      this.$http.get('/api/activities', { headers: headers, params: params }).then(response => {
        this.activities = response.body

        // Preselect all private/commute activities.
        this.selectedPrivate = []
        this.selectedCommute = []
        var that = this
        this.activities.forEach(function (activity) {
          if (activity.Private) {
            that.selectedPrivate.push(activity.Id)
          }
          if (activity.Commute) {
            that.selectedCommute.push(activity.Id)
          }
        })

        // Remember the original selections.
        this.fetchedPrivate = this.selectedPrivate.slice()
        this.fetchedCommute = this.selectedCommute.slice()

        this.selectCommutes()
      }, response => {
        this.error = response.statusText
      })
    },

    updateActivities: function () {
      if (this.modifiedPrivate.length === 0) {
        return
      }

      // List of activities with their new 'Private' flags.
      var updatePrivate = {}
      for (var i = 0; i < this.modifiedPrivate.length; i++) {
        var activityId = this.modifiedPrivate[i]
        updatePrivate[activityId] = this.selectedPrivate.includes(activityId)
      }

      // Disable button while request is being processed.
      this.modifiedPrivate = []

      var headers = auth.getAuthHeaders()
      var payload = { private: updatePrivate }
      this.$http.put('/api/activities_update', JSON.stringify(payload), { headers: headers }).then(response => {
        this.fetchActivities()
      }, response => {
        this.error = response.statusText
      })
    },

    selectCommutes: function () {
      if (!this.autoSelectCommutes) {
        return
      }

      if (this.startPlaceGeometry == null || this.endPlaceGeometry == null) {
        return
      }

      var placeStartCoords = {
        lat: this.startPlaceGeometry.location.lat(),
        lng: this.startPlaceGeometry.location.lng()
      }
      var placeEndCoords = {
        lat: this.endPlaceGeometry.location.lat(),
        lng: this.endPlaceGeometry.location.lng()
      }

      this.selectedCommute = []
      var that = this
      this.activities.forEach(function (activity) {
        var activityStartCoords = { lat: activity.Start_latlng[0], lng: activity.Start_latlng[1] }
        var activityEndCoords = { lat: activity.End_latlng[0], lng: activity.End_latlng[1] }

        var startToStart = haversineDistance(activityStartCoords, placeStartCoords)
        var startToEnd = haversineDistance(activityStartCoords, placeEndCoords)
        var endToStart = haversineDistance(activityEndCoords, placeStartCoords)
        var endToEnd = haversineDistance(activityEndCoords, placeEndCoords)

        if ((startToStart < that.proximityMeters && endToEnd < that.proximityMeters) ||
            (startToEnd < that.proximityMeters && endToStart < that.proximityMeters)) {
          that.selectedCommute.push(activity.Id)
        }
      })
    },

    setStartPlace: function (place) {
      this.startPlaceGeometry = place.geometry
      this.selectCommutes()
    },
    setEndPlace: function (place) {
      this.endPlaceGeometry = place.geometry
      this.selectCommutes()
    },

    getClassForActivity: function (id) {
      if (this.modifiedPrivate.includes(id) || this.modifiedCommute.includes(id)) {
        return 'is-selected'
      }

      return ''
    },
    isUpdateDisabled: function () {
      return this.modifiedPrivate.length === 0 && this.modifiedCommute.length === 0
    },

    formatDate: function (d) {
      var startDate = new Date(d)
      return dateFormat(startDate, 'd mmmm yyyy, HH:MM')
    },
    formatDistance: function (m) {
      return (m / 1000).toFixed(1)
    }
  },

  computed: {
    selectAllPrivate: {
      get: function () {
        return this.activities ? this.selectedPrivate.length === this.activities.length : false
      },
      set: function (value) {
        this.selectedPrivate = []
        var that = this
        if (value) {
          this.activities.forEach(function (activity) {
            that.selectedPrivate.push(activity.Id)
          })
        }
      }
    },
    selectAllCommutes: {
      get: function () {
        return this.activities ? this.selectedCommute.length === this.activities.length : false
      },
      set: function (value) {
        this.selectedCommute = []
        var that = this
        if (value) {
          this.activities.forEach(function (activity) {
            that.selectedCommute.push(activity.Id)
          })
        }
      }
    }
  }
}
</script>
