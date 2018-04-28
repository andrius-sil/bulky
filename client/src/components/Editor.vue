<template>
  <div>
    <datepicker v-model="startDate"></datepicker>
    <datepicker v-model="endDate"></datepicker>
    <table>
      <tr>
        <th><input type="checkbox" v-model="selectAll"></th>
        <th>Date</th>
        <th>Name</th>
        <th>Distance</th>
        <th>Commute</th>
        <th>Private</th>
      </tr>
      <tr v-for="activity in activities" :key="activity.Id">
        <td>
          <input type="checkbox" v-model="selected" :value="activity.Id" number>
        </td>
        <td>{{ formatDate(activity.Start_date_local) }}</td>
        <td>{{ activity.Name }}</td>
        <td>{{ formatDistance(activity.Distance) }}</td>
        <td>{{ formatBool(activity.Commute) }}</td>
        <td>{{ formatBool(activity.Private) }}</td>
      </tr>
    </table>
    <div v-if="activities.length == 0">No activities for selected dates.</div>
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

export default {
  data: function () {
    return {
      activities: [],
      selected: [],
      startDate: dateWeekAgo(),
      endDate: dateNow()
    }
  },

  components: {
    Datepicker
  },

  watch: {
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
    selectAll: {
      get: function () {
        return this.activities ? this.selected.length === this.activities.length : false
      },
      set: function (value) {
        var selected = []

        if (value) {
          this.activities.forEach(function (activity) {
            selected.push(activity.Id)
          })
        }

        this.selected = selected
      }
    }
  }
}
</script>
