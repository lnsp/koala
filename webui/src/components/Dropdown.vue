<template>
  <div class="relative w-full" @blur="active = false">
    <button
      @click="active = !active"
      class="border border-gray-300 w-full py-2 px-4 justify-between rounded flex items-center"
    >
      <span class="tracking-wide">{{ value }}</span>
      <span class="text-gray-500 text-xs">
        <span v-if="active">&#x25B2;</span>
        <span v-else>&#x25BC;</span>
      </span>
    </button>
    <div
      class="absolute z-50 mt-1 bg-white rounded overflow-hidden border border-gray-300 p-0 w-full"
      v-if="active"
    >
      <div
        v-for="opt in options"
        :key="opt"
        class="px-4 py-2 border-gray-300 tracking-wide hover:bg-gray-200 border-b last:border-b-0"
        @click="select(opt)"
      >{{ opt }}</div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Dropdown",
  props: ["value", "options"],
  data() {
    return {
      active: false
    };
  },
  methods: {
    select(option) {
      this.value = option;
      this.active = false;
      this.$emit("change");
    }
  }
};
</script>