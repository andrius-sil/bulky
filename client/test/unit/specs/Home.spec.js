import Vue from 'vue'
import Home from '@/components/Home'

import { shallow, createLocalVue } from '@vue/test-utils'
import VueRouter from 'vue-router'

describe('Home.vue', () => {
  it('should render correct contents', () => {
    const localVue = createLocalVue()
    localVue.use(VueRouter)
    const router = new VueRouter()
    const wrapper = shallow(Home, {
        localVue,
        router
    })
    expect(wrapper.vm.$el.querySelector('h1').textContent)
      .toEqual('App needs Strava permissions to see private activities and modify all activities.')
  })
})
