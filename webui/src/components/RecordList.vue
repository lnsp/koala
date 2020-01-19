<template>
  <div>
    <div
      class="hidden sm:flex flex-row items-center uppercase text-gray-500 font-bold tracking-wider text-xs mb-4"
    >
      <div class="w-1/5 px-4">Type</div>
      <div class="w-1/3 px-8">Domain</div>
      <div class="w-1/3 px-8">Data</div>
    </div>
    <div>
      <div
        v-for="(record, index) in filtered"
        :key="index"
        class="flex flex-row items-center py-2 flex-wrap odd:bg-gray-100"
      >
        <div class="w-1/2 sm:w-1/5 px-3 mb-3 sm:mb-0">
          <dropdown v-model="record.type" :options="['A', 'CNAME']" />
        </div>
        <div class="w-1/2 sm:w-1/3 px-3 mb-3 sm:mb-0">
          <input
            type="text"
            placeholder="Domain"
            v-model="record.domain"
            class="w-full py-2 px-4 bg-transparent focus:outline-none focus:shadow-outline appearance-none leading-normal border border-transparent hover:border-gray-300 focus:bg-gray-100"
          />
        </div>
        <div class="flex-auto px-3">
          <input
            type="text"
            placeholder="Data"
            v-model="record.data"
            class="w-full py-2 px-4 bg-transparent focus:outline-none focus:shadow-outline appearance-none leading-normal border border-transparent hover:border-gray-300 focus:bg-gray-100"
          />
        </div>
        <div class="px-6">
          <button class="text-gray-500 focus:outline-none focus:shadow-outline hover:text-indigo-600 text-2xl" @click="discard(record)">&#215;</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Dropdown from "./Dropdown";
import TextInput from "./TextInput";

export default {
  name: "RecordList",
  props: ["records", "filter"],
  components: {
    dropdown: Dropdown,
    textInput: TextInput
  },
  computed: {
    filtered: function() {
      if (this.filter === '') return this.records;
      return this.records.filter(record => {
        let f = this.filter.toLowerCase()
        return record.type.toLowerCase().indexOf(f) >= 0 || record.domain.toLowerCase().indexOf(f) >= 0 || record.data.toLowerCase().indexOf(f) >= 0
      })
    }
  },
  methods: {
    discard(record) {
      this.records.splice(this.records.indexOf(record), 1);
    },
  }
};
</script>