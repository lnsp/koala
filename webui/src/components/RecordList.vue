<template>
  <div>
    <div
      class="hidden sm:flex flex-row items-center uppercase text-gray-500 font-bold tracking-wider text-xs p-3 sm:px-6"
    >
      <div class="w-1/5">Type</div>
      <div class="w-1/4 mx-3 px-4">Domain</div>
      <div class="w-1/3 mx-3 px-4">Data</div>
    </div>
    <div class="mb-4">
      <div
        v-for="(record, index) in filtered"
        :key="index"
        class="relative flex flex-row sm:items-center p-3 sm:px-6 flex-wrap sm:odd:bg-gray-100 border sm:border-transparent mb-3 mx-3 sm:m-0"
      >
        <div class="w-full sm:w-1/5 mb-3 sm:mb-0">
          <label class="block sm:hidden uppercase text-gray-500 font-bold tracking-wider text-xs px-4 py-2">Type</label>
          <dropdown v-model="record.type" :options="['A', 'CNAME', 'AAAA']" />
        </div>
        <div class="w-full sm:w-1/4 mx-0 sm:mx-3 mb-3 sm:mb-0">
          <label class="block sm:hidden uppercase text-gray-500 font-bold tracking-wider text-xs px-4 py-2">Domain</label>
          <input
            type="text"
            placeholder="Domain"
            v-model="record.domain"
            class="w-full py-2 px-4 bg-transparent focus:outline-none focus:shadow-outline appearance-none leading-normal border border-gray-300 sm:border-transparent hover:border-gray-300 focus:bg-gray-100"
          />
        </div>
        <div class="flex-grow mx-0 sm:mx-3">
          <label class="block sm:hidden uppercase text-gray-500 font-bold tracking-wider text-xs px-4 py-2">Data</label>
          <input
            type="text"
            placeholder="Data"
            v-model="record.data"
            class="w-full py-2 px-4 bg-transparent focus:outline-none focus:shadow-outline appearance-none leading-normal border border-gray-300 sm:border-transparent hover:border-gray-300 focus:bg-gray-100"
          />
        </div>
        <div class="absolute sm:relative pl-3 right-0 mr-3 -mt-2 sm:mr-0 sm:-mt-0">
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