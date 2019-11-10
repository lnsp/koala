<template>
  <div id="app">
    <router-view />
  </div>
</template>

<style>
#app {
  font-family: -apple-system, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
.btn-koala {
  color: #fff !important;
  background-color: rgba(117,72,232,1) !important;
  border-color: rgba(117,72,232,1) !important;
}
html {
  height: 100%;
}
body {
  background-color: none;
}
</style>

<style lang="scss">
@import '../node_modules/bootstrap/scss/bootstrap.scss';
</style>