<template>
  <div id="app">
    <ControlPanel/>
  </div>
</template>

<script>
import ControlPanel from './components/ControlPanel.vue'

export default {
  name: 'app',
  components: {
    ControlPanel
  }
}
</script>

<style>
#app {
  font-family: Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>

<style lang="scss">
@import '../node_modules/bootstrap/scss/bootstrap.scss';
</style>