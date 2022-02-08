import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import ErrorPage from '@/views/ErrorPage.vue'

Vue.use(Router)

export default new Router({
    mode: 'history',
    base: process.env.BASE_URL,
    routes: [
        {
            path: '/',
            name: 'home',
            component: Home
        },
        {
            path: '/login',
            name: 'login',
            component: () =>
                import(/* webpackChunkName: "login" */ '@/views/Login.vue')
        },
        {
            path: '/forget',
            name: 'forget',
            component: () =>
                import(/* webpackChunkName: "forget" */ '@/views/Forget.vue')
        },
        {
            path: '/reset',
            name: 'reset',
            component: () =>
                import(/* webpackChunkName: "reset" */ '@/views/Reset.vue')
        },
        {
            path: '/bulletin',
            name: 'bulletin',
            component: () =>
                import(
                    /* webpackChunkName: "bulletin" */ '@/views/Bulletin.vue'
                )
        },
        {
            path: '/rank/:page',
            name: 'rank',
            component: () =>
                import(/* webpackChunkName: "rank" */ '@/views/Rank.vue')
        },
        {
            path: '/profile/:user_id',
            name: 'profile',
            component: () =>
                import(/* webpackChunkName: "profile" */ '@/views/Profile.vue'),
            meta: {
                auth: 'member'
            }
        },
        {
            path: '/challenge',
            name: 'challenge',
            component: () =>
                import(
                    /* webpackChunkName: "challenge" */ '@/views/Challenge.vue'
                ),
            meta: {
                auth: 'member'
            }
        },
        {
            path: '/admin/user',
            name: 'admin-user',
            component: () =>
                import(
                    /* webpackChunkName: "admin-challenge" */ '@/views/admin/User.vue'
                ),
            meta: {
                auth: 'admin'
            }
        },
        {
            path: '/admin/challenge',
            name: 'admin-challenge',
            component: () =>
                import(
                    /* webpackChunkName: "admin-challenge" */ '@/views/admin/Challenge.vue'
                ),
            meta: {
                auth: 'team'
            }
        },
        {
            path: '*',
            name: 'error',
            component: ErrorPage
        }
        // {
        //   path: "/about",
        //   name: "about",
        //   // route level code-splitting
        //   // this generates a separate chunk (about.[hash].js) for this route
        //   // which is lazy-loaded when the route is visited.
        //   component: () =>
        //     import(/* webpackChunkName: "about" */ "./views/About.vue")
        // }
    ]
})
