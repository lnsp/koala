<template>
  <div id="app" class="h-screen w-screen">
    <router-view />
  </div>
</template>

<style>
@tailwind base;
@tailwind components;
@tailwind utilities;

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