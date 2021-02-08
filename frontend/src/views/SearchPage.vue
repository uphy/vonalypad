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
    <v-dialog
      v-model="dialog"
      fullscreen
      hide-overlay
      transition="dialog-bottom-transition"
      v-if="recipe"
    >
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          color="red lighten-2"
          dark
          v-bind="attrs"
          v-on="on"
        >
          Click Me
        </v-btn>
      </template>

      <v-card>
        <v-card-title class="headline grey lighten-2">
          {{ recipe.title}}
          <v-spacer></v-spacer>
          <v-btn @click="close">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>

        <v-card-text>
          {{ recipe.description }}
        </v-card-text>

        <v-divider></v-divider>
        <div style="padding: 1rem">
          <ul>
            <li
              v-for="(ingredient, i) in recipe.ingredients"
              :key="i"
            >{{ ingredient. name}}: {{ ingredient.quantity }}</li>
          </ul>
        </div>

        <v-divider></v-divider>
        <div style="padding: 1rem">
          <ol>
            <li
              v-for="(step, i) in recipe.steps"
              :key="i"
            >{{ step.step}}
              <v-img
                v-if="step.image"
                :src="step.image"
                max-width="300px"
              />
            </li>
          </ol>
        </div>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import axios from "axios";
import NoSleep from "nosleep.js";

const noSleep = new NoSleep();

export default {
  data: () => {
    return {
      query: "",
      recipes: [],
      recipe: null,
      dialog: false,
    };
  },
  mounted() {
    this.random();
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
      this.recipes = this.filterRecipe(resp.data);
    },
    filterRecipe(recipe) {
      recipe.forEach((recipe) => {
        recipe.show = false;
      });
      return recipe;
    },
    show(recipe) {
      this.recipe = recipe;
      this.dialog = true;
      noSleep.enable();
    },
    close() {
      noSleep.disable();
      this.dialog = false;
      this.recipe = null;
    },
  },
};
</script>


<style scoped>
</style>