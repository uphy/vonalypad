<template>
  <v-combobox
    v-model="tags"
    :items="allTags"
    chips
    clearable
    label="Tags"
    multiple
    solo
  >
    <template v-slot:selection="{ attrs, item, select, selected }">
      <v-chip
        v-bind="attrs"
        :input-value="selected"
        close
        @click="select"
        @click:close="removeTag(item)"
      >
        {{ item }}
      </v-chip>
    </template>
  </v-combobox>
</template>

<script>
import axios from "axios";

export default {
  props: {
    recipeId: String,
  },
  data: () => {
    return {
      tags: [],
      allTags: [],
    };
  },
  mounted() {
    this.fetchAllTags();
    this.fetchSelectedTags();
  },
  watch: {
    async tags(newTags, oldTags) {
      console.log(newTags);
      if (newTags.length !== oldTags.length) {
        console.log(newTags);
      }
      //await this.fetchSelectedTags();
    },
  },
  methods: {
    async fetchAllTags() {
      const resp = await axios.get(`./api/tags/`);
      this.allTags = resp.data.map((tag) => tag.name);
    },
    async fetchSelectedTags() {
      const resp = await axios.get(`./api/recipes/${this.recipeId}/tags`);
      this.tags = resp.data.map((tag) => tag.name);
    },
    async updateTag() {
      const resp = await axios.put(
        `./api/recipes/${this.recipeId}/tags`,
        this.tags
      );
    },
    removeTag(item) {
      console.log(`remove tag: ${item}`);
    },
  },
};
</script>