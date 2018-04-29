<template>
  <div>
    <datepicker v-model="startDate"></datepicker>
    <datepicker v-model="endDate"></datepicker>
    <table border="1">
      <tr>
        <th>Date</th>
        <th>Name</th>
        <th>Distance</th>
        <th>Commute</th>
        <th>
          <input type="checkbox" v-model="selectAllPrivate">
          Private
        </th>
      </tr>
      <tr v-for="activity in activities" :key="activity.Id">
        <td>{{ formatDate(activity.Start_date_local) }}</td>
        <td>{{ activity.Name }}</td>
        <td>{{ formatDistance(activity.Distance) }}</td>
        <td>{{ formatBool(activity.Commute) }}</td>
        <td :class="modifiedPrivate.includes(activity.Id) ? 'private' : ''">
          <input type="checkbox" v-model="selectedPrivate" :value="activity.Id" number>
        </td>
      </tr>
    </table>
    <div v-if="activities.length === 0">No activities for selected dates.</div>
    <button :disabled="modifiedPrivate.length === 0" @click="updateActivities()">Update Activities</button>
  </div>
</template>

<script>
import dateFormat from 'dateformat'
import Datepicker from 'vuejs-datepicker'

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

      startDate: dateWeekAgo(),
      endDate: dateNow()
    }
  },

  components: {
    Datepicker
  },

  watch: {
    selectedPrivate: function (val) {
      this.modifiedPrivate = arrayDiff(this.selectedPrivate, this.fetchedPrivate)
    },
    startDate: function (val) {
      this.fetchActivities()
    },
    endDate: function (val) {
      this.fetchActivities()
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

        // Preselect all private activities.
        this.selectedPrivate = []
        var that = this
        this.activities.forEach(function (activity) {
          if (activity.Private) {
            that.selectedPrivate.push(activity.Id)
          }
        })

        // Remember the original selection.
        this.fetchedPrivate = this.selectedPrivate
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

      // This is to disable button while request is being processed.
      this.modifiedPrivate = []

      var headers = auth.getAuthHeaders()
      var payload = { private: updatePrivate }
      this.$http.put('/api/activities_update', JSON.stringify(payload), { headers: headers }).then(response => {
        this.fetchActivities()
      }, response => {
        this.error = response.statusText
      })
    },

    formatDate: function (d) {
      var startDate = new Date(d)
      return dateFormat(startDate, 'd mmmm yyyy, HH:MM')
    },

    formatDistance: function (m) {
      return (m / 1000).toFixed(1)
    },

    formatBool: function (c) {
      if (c) {
        return 'yes'
      } else {
        return 'no'
      }
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
    }
  }
}
</script>

<style scoped>
td.private {
  background-color: #369;
}
</style>
