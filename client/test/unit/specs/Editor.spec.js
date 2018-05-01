import Vue from 'vue'
import Editor from '@/components/Editor'

import { shallow, createLocalVue } from '@vue/test-utils'
import VueRouter from 'vue-router'
import VueResource from 'vue-resource'

describe('Editor.vue', () => {
  it('should render correct contents', () => {
    const localVue = createLocalVue()
    localVue.use(VueRouter)
    localVue.use(VueResource)
    const router = new VueRouter()
    const wrapper = shallow(Editor, {
        localVue,
        router
    })
    expect(wrapper.vm.$el.querySelector('.no-activities').textContent)
      .toEqual('No activities for selected dates.')
  })
})
