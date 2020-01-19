<template>
  <div class="h-screen p-10 bg-gray-900">
    <div class="bg-white rounded-lg max-w-lg px-6 py-10 mx-auto shadow-lg">
      <div class="mb-10 text-center">
        <brand/>
        <div class="text-gray-600">Please authenticate to continue.</div>
      </div>
      <form @submit.prevent="finish()">
        <div class="sm:flex flex-wrap">
          <input type="text" placeholder="Token"
            class="bg-gray-200 font-mono focus:outline-none focus:shadow-outline rounded py-2 px-4 appearance-none w-full sm:w-auto sm:flex-grow mb-3 sm:mb-0 sm:mr-3"
            v-model="token" />
          <button type="submit" class="bg-indigo-600 font-medium text-white hover:bg-indigo-700 focus:outline-none focus:shadow-outline rounded py-2 px-4 shadow">Sign in</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import Brand from '../components/Brand'

export default {
  name: "Authentication",
  data() {
    return {
      token: localStorage.token
    };
  },
  components: {
    'brand': Brand
  },
  methods: {
    finish: function() {
      localStorage.token = this.token;
      this.$router.push("/");
    }
  }
};
</script>