import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import { createApolloProvider } from './plugins/vue-apollo'
import '@/main.scss'

Vue.config.productionTip = false

router.beforeEach((to, from, next) => {
    let validator =
        typeof to.meta!.auth != 'string' ||
        (to.meta!.auth == 'member' && !!sessionStorage.getItem('user_id')) ||
        (to.meta!.auth == 'team' &&
            sessionStorage.getItem('role') != 'member') ||
        (to.meta!.auth == 'admin' && sessionStorage.getItem('role') == 'admin')
    let result = validator
        ? {}
        : {
              name: 'home'
          }
    next(result)
})

new Vue({
    vuetify,
    router,
    store,
    apolloProvider: createApolloProvider(),
    render: (h) => h(App)
}).$mount('#app')
