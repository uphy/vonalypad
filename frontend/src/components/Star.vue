<template>
  <v-icon
    right
    dark
    @click="toggle"
    :color="checked?'yellow':'white'"
  >
    mdi-star
  </v-icon>
</template>

<script>
import axios from "axios";

const tagName = "Starred";
export default {
  props: {
    recipeId: String,
    starred: Boolean,
  },
  data: () => {
    return {
      checked: false,
    };
  },
  mounted() {
    if (this.starred === undefined) {
      this.fetchStarred();
    } else {
      this.checked = this.starred;
    }
  },
  methods: {
    async toggle() {
      if (this.checked) {
        await axios.post(`./api/recipes/${this.recipeId}/remove-tag`, {
          tag: tagName,
        });
      } else {
        await axios.post(`./api/recipes/${this.recipeId}/add-tag`, {
          tag: tagName,
        });
      }
      await this.fetchStarred();
    },
    async fetchStarred() {
      const resp = await axios.get(`./api/recipes/${this.recipeId}/tags`);
      for (const tag of resp.data) {
        if (tag.name === tagName) {
          this.checked = true;
          return;
        }
      }
      this.checked = false;
    },
  },
};
</script>