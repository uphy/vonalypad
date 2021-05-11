<template>
  <v-container fluid>
    <v-card>
      <v-card-title>Note</v-card-title>
      <v-card-text>
        <v-textarea
          auto-grow
          v-model="note"
          rows="1"
        ></v-textarea>
      </v-card-text>
      <v-card-actions>
        <v-btn
          depressed
          @click="update"
        >
          Save
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script>
import axios from "axios";

export default {
  props: {
    recipeId: String,
  },
  data: () => {
    return {
      note: "",
    };
  },
  async mounted() {
    await this.fetch();
  },
  watch: {
    async recipeId() {
      await this.fetch();
    },
  },
  methods: {
    async fetch() {
      const resp = await axios.get(`./api/recipes/${this.recipeId}/note`);
      this.note = resp.data;
    },
    async update() {
      await axios.put(`./api/recipes/${this.recipeId}/note`, {
        note: this.note,
      });
    },
  },
};
</script>