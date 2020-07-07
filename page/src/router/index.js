import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '@/views/Home.vue'
import Console from '@/views/Console.vue'
import Desktop from '@/views/Desktop.vue'
import NotFound from '@/views/NotFound.vue'

Vue.use(VueRouter)

export default new VueRouter({
  mode: 'hash',
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home
    },
    {
      path: '/console',
      name: 'Console',
      component: Console
    },
    {
      path: '/desktop',
      name: 'Desktop',
      component: Desktop
    },
    {
      path: '*',
      name: NotFound,
      component: NotFound
    }
  ]
})