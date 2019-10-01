import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import { createApolloProvider } from './plugins/vue-apollo'
import '@/main.scss'

Vue.config.productionTip = false

new Vue({
    vuetify,
    router,
    store,
    apolloProvider: createApolloProvider(),
    render: (h) => h(App)
}).$mount('#app')
