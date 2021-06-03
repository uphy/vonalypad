<template>
  <v-form v-model="valid">
    <v-container>
      <v-row>
        <v-col cols="12" md="4">
          <v-text-field
            v-model="url"
            label="Cookpad URL"
            required
          ></v-text-field>
        </v-col>
      </v-row>
      <v-btn :disabled="!valid" color="success" @click="importUrl">
        Import URL
      </v-btn>
    </v-container>
  </v-form>
</template>

<script>
import axios from "axios";

export default {
  data: () => {
    return {
      url: "",
      valid: false,
    };
  },
  methods: {
    async importUrl() {
      const recipe = await axios.post("./api/recipes/import", {
        url: this.url,
      });
      this.show(recipe.data);
    },
    show(recipe) {
      this.$router.push(`/recipes/${recipe.id}`);
    },
  },
};
</script>
