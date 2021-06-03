<template>
  <div>
    <v-container>
      <v-responsive max-width="300">
        <v-text-field
          dense
          flat
          hide-details
          rounded
          solo
          append-icon="mdi-magnify"
          v-model="query"
          @keyup.enter="search"
          @click="search"
        ></v-text-field>
      </v-responsive>
    </v-container>
    <v-row>
      <div v-if="loading">
        <v-progress-circular
          indeterminate
          color="primary"
        ></v-progress-circular>
      </div>
      <v-col v-else v-for="r in recipes" :key="r.id" cols="6" md="3">
        <v-card @click="show(r)">
          <v-img :src="r.image" width="100%" aspect-ratio="1"> </v-img>
          <v-card-title style="font-size: 0.8rem; line-height: 1rem">
            {{ r.title }}
            <div>🍳{{ r.tsukurepo }}<span v-if="r.video"> 🎥</span></div>
          </v-card-title>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import axios from "axios";

let randomRecipes = null;

export default {
  props: {
    q: String,
  },
  data: () => {
    return {
      query: "",
      recipes: [],
      recipe: null,
      loading: false,
    };
  },
  mounted() {
    this.updateByProps();
  },
  watch: {
    q() {
      this.updateByProps();
    },
  },
  methods: {
    async updateByProps() {
      if (this.q === undefined || this.q === "") {
        if (randomRecipes === null) {
          randomRecipes = await this.random();
        }
        this.query = "";
        this.recipes = randomRecipes;
      } else {
        this.query = this.q;
        await this.search();
      }
    },
    async search() {
      if (this.query.startsWith("tag:")) {
        const tag = this.query.substring(this.query.indexOf(":") + 1);
        this.loading = true;
        try {
          const resp = await axios.get(`./api/tags/${tag}`, {
            params: {
              q: this.query,
            },
          });
          this.recipes = this.filterRecipe(resp.data);
        } finally {
          this.loading = false;
        }
        return;
      }

      this.loading = true;
      try {
        const resp = await axios.get("./api/search?", {
          params: {
            q: this.query,
          },
        });
        this.recipes = this.filterRecipe(resp.data);
      } finally {
        this.loading = false;
      }
    },
    async random() {
      this.loading = true;
      try {
        const resp = await axios.get("./api/random");
        return this.filterRecipe(resp.data);
      } finally {
        this.loading = false;
      }
    },
    filterRecipe(recipe) {
      recipe.forEach((recipe) => {
        recipe.show = false;
      });
      return recipe;
    },
    show(recipe) {
      this.$router.push(`/recipes/${recipe.id}`);
    },
  },
};
</script>

<style scoped></style>
