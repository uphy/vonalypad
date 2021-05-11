<template>
  <v-card :loading="loading">
    <v-card-title>Note</v-card-title>
    <v-card-text>
      <v-textarea
        auto-grow
        v-model="note"
        rows="1"
        v-if="editorMode"
        @blur="blurred"
      ></v-textarea>
      <div
        v-html="compiledMarkdown"
        @click="editorMode=true"
        v-else
        style="min-height:1rem"
      />
    </v-card-text>
    <v-card-actions>
      <v-btn
        depressed
        @click="update"
        v-if="editorMode"
      >
        Save
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import axios from "axios";
import marked from "marked";

export default {
  props: {
    recipeId: String,
  },
  data: () => {
    return {
      note: "",
      loading: false,
      editorMode: false,
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
  computed: {
    compiledMarkdown() {
      return marked(this.note);
    },
  },
  methods: {
    blurred() {
      window.setTimeout(() => {
        this.editorMode = false;
      }, 100);
    },
    async fetch() {
      this.loading = true;
      try {
        const resp = await axios.get(`./api/recipes/${this.recipeId}/note`);
        this.note = resp.data;
      } finally {
        this.loading = false;
      }
      if (this.note.length === 0) {
        this.editorMode = true;
      }
    },
    async update() {
      this.loading = true;
      try {
        await axios.put(`./api/recipes/${this.recipeId}/note`, {
          note: this.note,
        });
      } finally {
        this.loading = false;
      }
      this.editorMode = false;
    },
  },
};
</script>