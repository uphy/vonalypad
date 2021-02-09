import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import SearchPage from "@/views/SearchPage.vue";
import RecipePage from "@/views/RecipePage.vue";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "SearchPage",
    component: SearchPage,
  },
  {
    path: "/recipes/:recipeId",
    name: "RecipePage",
    component: RecipePage,
  }
];

const router = new VueRouter({
  routes
});

export default router;
