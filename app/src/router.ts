import Vue from "vue";
import Router from "vue-router";
import Home from "./views/Home.vue";
import ErrorPage from "@/views/ErrorPage.vue";

Vue.use(Router);

export default new Router({
  mode: "history",
  base: process.env.BASE_URL,
  routes: [
    {
      path: "/",
      name: "home",
      component: Home
    },
    {
      path:'/login',
      name: "login",
      component : ()=>import(/* webpackChunkName: "login" */ "@/views/Login.vue")
    },
    {
      path:'/bulletin',
      name: "bulletin",
      component : ()=>import(/* webpackChunkName: "bulletin" */ "@/views/Bulletin.vue")
    },
    {
      path:'/rank/:page',
      name: "rank",
      component : ()=>import(/* webpackChunkName: "rank" */ "@/views/Rank.vue")
    },
    {
      path:'/challenge',
      name: "challenge",
      component : ()=>import(/* webpackChunkName: "challenge" */ "@/views/Challenge.vue")
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
});
