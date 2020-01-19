<template>
  <div>
    <div
      class="hidden sm:flex flex-row items-center uppercase text-gray-500 font-bold tracking-wider text-xs mb-4"
    >
      <div class="w-1/5 px-3">Type</div>
      <div class="w-1/3 px-3">Domain</div>
      <div class="w-1/3 px-3">Data</div>
    </div>
    <div v-for="(record, index) in records" :key="index">
      <div class="flex flex-row items-center mb-10 sm:mb-4 flex-wrap">
        <div class="w-1/2 sm:w-1/5 px-3 mb-3 sm:mb-0">
          <dropdown v-model="record.type" :options="['A', 'CNAME']" />
        </div>
        <div class="w-1/2 sm:w-1/3 px-3 mb-3 sm:mb-0">
          <input
            type="text"
            placeholder="Domain"
            v-model="record.domain"
            class="w-full py-2 px-4 bg-white focus:outline-none focus:shadow-outline border appearance-none leading-normal border-gray-200 bg-gray-100"
          />
        </div>
        <div class="flex-auto px-3">
          <input
            type="text"
            placeholder="Data"
            v-model="record.data"
            class="w-full py-2 px-4 bg-white focus:outline-none focus:shadow-outline border appearance-none leading-normal border-gray-200 bg-gray-100"
          />
        </div>
        <div class="px-3">
          <button class="text-gray-500 text-lg" @click="discard(index)">&#215;</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Dropdown from "./Dropdown";

export default {
  name: "RecordList",
  props: ["records"],
  components: {
    dropdown: Dropdown
  },
  methods: {
    discard(index) {
      this.records.splice(index, 1);
    }
  }
};
</script>