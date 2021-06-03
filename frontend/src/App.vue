<template>
  <v-app>
    <v-app-bar app>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title @click="$router.push('/')" style="cursor:pointer;"
        >Vonalypad</v-toolbar-title
      >
    </v-app-bar>

    <v-navigation-drawer v-model="drawer" fixed temporary>
      <v-list>
        <v-list-item
          v-for="link in links"
          :key="link.text"
          link
          :to="link.link"
        >
          <v-list-item-icon>
            <v-icon>{{ link.icon }}</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ link.text }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-main class="grey lighten-2">
      <v-container>
        <router-view />
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import axios from "axios";

export default Vue.extend({
  name: "App",
  data: () => ({
    drawer: null,
    links: [
      { icon: "mdi-magnify", text: "Search", link: "/" },      
    ],
  }),
  async mounted() {
    const links = [];
    links.push({ icon: "mdi-magnify", text: "Search", link: "/" });
    links.push({ icon: "mdi-magnify", text: "Import", link: "/import" });
    const tags = await axios.get("./api/tags/");
    for (const tag of tags.data) {
      links.push({ icon: "mdi-tag", text: tag.name, link: `/tags/${tag.id}` });
    }
    this.links = links;
  },
});
</script>
