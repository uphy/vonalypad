import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import SearchPage from "@/views/SearchPage.vue";
import RecipePage from "@/views/RecipePage.vue";
import ImportPage from "@/views/ImportPage.vue";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "SearchPage",
    component: SearchPage,
    props: route => ({ q: '' })
  },
  {
    path: "/recipes/:recipeId",
    name: "RecipePage",
    component: RecipePage,
  },
  {
    path: "/import",
    name: "ImportPage",
    component: ImportPage,
  },
  {
    path: "/tags/:q",
    name: "SearchPage",
    component: SearchPage,
    props: route => ({ q: `tag:${route.params.q}` })
  }
];

const router = new VueRouter({
  routes
});

export default router;
