<template>
  <div>
    <input
      type="text"
      v-model="state.query"
    >
    <input
      type="submit"
      value="Search"
      @click="search"
    >
    <div class="recipes">
      <div
        class="recipe"
        v-for="r in state.recipes"
        :key="r.id"
      >
        <a
          :href="'https://cookpad.com/recipe/'+r.id"
          target="_blank"
        ><img
            class="image"
            :src="r.image"
          ></a>
        <div class="title"><span>{{r.title }}</span></div>
        <span class="tsukurepo">🍳{{ r.tsukurepo }}<span v-if="r.video"> 🎥</span></span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.recipes {
  font-size: 0.8rem;
}
.recipe {
  margin: 0.3rem;
  display: inline-block;
  position: relative;
}
.description {
  font-size: 0.6rem;
  color: gray;
}
.image {
  width: 256px;
  height: 256px;
  display: block;
}
.title {
  position: absolute;
  bottom: 0;
  width: 256px;
  background: rgba(0, 0, 0, 0.5);
  color: white;
}
.title > span {
  margin: 0.2rem 0.5rem;
}
.tsukurepo {
  position: absolute;
  top: 0;
  right: 0;
  max-width: 256px;
  background: rgba(226, 147, 0, 0.5);
  color: white;
  padding: 0.2rem 0.5rem;
}
</style>

<script>
import axios from "axios";
import { reactive } from "vue";

export default {
  name: "Search",
  data() {
    const state = reactive({
      recipes: [],
      query: "鶏肉",
    });
    return {
      state,
      search: async () => {
        const resp = await axios.get("./api/search?", {
          params: {
            q: state.query,
          },
        });
        state.recipes = resp.data;
      },
    };
  },
};
</script>
