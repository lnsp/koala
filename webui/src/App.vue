<template>
  <div id="app" class="min-h-screen">
    <ControlPanel />
  </div>
</template>

<script>
import ControlPanel from './components/ControlPanel'

export default {
  name: 'App',
  components: { ControlPanel }
}
</script>


<style>
@tailwind base;
@tailwind components;
@tailwind utilities;

html, body {
  @apply bg-gray-900;
}

@font-face {
  font-family: 'Inter';
  src: url('./assets/fonts/Inter-Regular.ttf');
  font-weight: 400;
}

@font-face {
  font-family: 'Inter';
  src: url('./assets/fonts/Inter-Medium.ttf');
  font-weight: 500;
}

@font-face {
  font-family: 'Inter';
  src: url('./assets/fonts/Inter-SemiBold.ttf');
  font-weight: 600;
}
</style>