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
      <v-col
        v-for="r in recipes"
        :key="r.id"
        cols="6"
        md="3"
      >
        <v-card @click="show(r)">
          <v-img
            :src="r.image"
            width="100%"
            aspect-ratio="1"
          >
          </v-img>
          <v-card-title style="font-size: 0.8rem; line-height: 1rem">
            {{r.title}}
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
  data: () => {
    return {
      query: "",
      recipes: [],
      recipe: null,
    };
  },
  async mounted() {
    if (randomRecipes === null) {
      randomRecipes = await this.random();
    }
    this.recipes = randomRecipes;
  },
  methods: {
    async search() {
      const resp = await axios.get("./api/search?", {
        params: {
          q: this.query,
        },
      });
      this.recipes = this.filterRecipe(resp.data);
    },
    async random() {
      const resp = await axios.get("./api/random");
      return this.filterRecipe(resp.data);
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


<style scoped>
</style>