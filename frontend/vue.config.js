module.exports = {
  transpileDependencies: ["vuetify"],
  pages: {
    index: {
      entry: "src/main.ts",
      title: "Vonalypad",
    }
  },
  devServer: {
    port: 8081,
    host: '0.0.0.0',
    disableHostCheck: true,
    proxy: 'http://localhost:8080'
  },
  publicPath: "./"
};
