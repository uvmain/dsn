<script setup lang="ts">
</script>

<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <NavBar />
    <main class="py-8">
      <RouterView />
    </main>
    <NotificationContainer />
  </div>
</template>

<style>
html, body, #app {
  height: 100%;
  margin: 0;
  padding: 0;
}

#app {
  font-family: 'Inter', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
