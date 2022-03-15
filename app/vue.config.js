const VuetifyLoaderPlugin = require('vuetify-loader/lib/plugin');

module.exports = {
  lintOnSave: false,
  devServer: {
    proxy: {
      '/graphql': {
        target: "http://localhost:8888"
      },
      '/writeup': {
        target: "http://localhost:8888"
      },
      '/allwp': {
        target: "http://localhost:8888"
      },
      '/wp': {
        target: "http://localhost:8888"
      }
    }
  },
  configureWebpack: {
    plugins: [
      new VuetifyLoaderPlugin()
    ],
  },
};
