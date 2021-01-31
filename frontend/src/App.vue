<template>
  <div>
    <h1>Vonalypad</h1>
    <Search />
  </div>
</template>

<script>
import Search from "./components/Search.vue";

export default {
  name: "App",
  components: {
    Search,
  },
};
</script>
