import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import SearchPage from "@/views/SearchPage.vue";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "SearchPage",
    component: SearchPage,
  }
];

const router = new VueRouter({
  routes
});

export default router;
