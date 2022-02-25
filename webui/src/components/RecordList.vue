<template>
  <div>
    <div class="flex flex-row p-3 sm:px-6 sm:pt-6">
      <button class="text-white px-4 py-2 mr-3 rounded bg-indigo-600 hover:bg-indigo-700 shadow"
              @click="addRecord">
        <span class="inline sm:hidden">Add</span>
        <span class="hidden sm:inline">Add Record</span>
      </button>
      <textInput class="flex-grow"
                 v-model="filter"
                 placeholder="Search for record" />
    </div>
    <div class="hidden sm:flex flex-row items-center uppercase text-gray-500 font-bold tracking-wider text-xs p-3 sm:px-6">
      <div class="w-1/5">Type</div>
      <div class="w-1/4 mx-3 px-4">Domain</div>
      <div class="w-1/3 mx-3 px-4">Data</div>
    </div>
    <div class="mb-4">
      <div v-for="record in filtered"
           :key="record.type"
           class="relative flex flex-row sm:items-center p-3 sm:px-6 xs:flex-wrap sm:odd:bg-gray-100 rounded border sm:border-transparent mb-3 mx-3 sm:m-0">
        <div class="w-full sm:w-1/5 mb-3 sm:mb-0">
          <label class="block sm:hidden uppercase text-gray-500 font-bold tracking-wider text-xs px-4 py-2">Type</label>
          <record-selector v-model="record.type" />
        </div>
        <div class="w-full sm:w-1/4 mx-0 sm:mx-3 mb-3 sm:mb-0">
          <label class="block sm:hidden uppercase text-gray-500 font-bold tracking-wider text-xs px-4 py-2">Domain</label>
          <input type="text"
                 placeholder="Domain"
                 v-model="record.domain"
                 class="w-full py-2 px-4 bg-transparent focus:outline-none focus:shadow-outline appearance-none leading-normal border border-gray-300 sm:border-transparent hover:border-gray-300 focus:bg-gray-100" />
        </div>
        <div class="flex-grow mx-0 sm:mx-3">
          <label class="block sm:hidden uppercase text-gray-500 font-bold tracking-wider text-xs px-4 py-2">Data</label>
          <input type="text"
                 placeholder="Data"
                 v-model="record.data"
                 class="w-full py-2 px-4 bg-transparent focus:outline-none focus:shadow-outline appearance-none leading-normal border border-gray-300 sm:border-transparent hover:border-gray-300 focus:bg-gray-100" />
        </div>
        <div class="absolute sm:relative pl-3 right-0 mr-3 -mt-2 sm:mr-0 sm:-mt-0">
          <button class="text-gray-500 focus:outline-none focus:shadow-outline hover:text-indigo-600 text-2xl"
                  @click="discardRecord(record)">&#215;</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import RecordSelector from "./RecordSelector.vue";
import TextInput from "../components/TextInput.vue";

export default {
  props: ['zone'],
  components: { RecordSelector, TextInput },
  data () {
    return {
      filter: '',
    }
  },
  methods: {
    addRecord () {
      this.$store.commit('addRecord', { zone: this.zone, record: { type: 'A', domain: '', data: '' }})
    },
    discardRecord (record) {
      this.$store.commit('dropRecord', { zone: this.zone, record: record })
    },
  },
  computed: {
    records () {
      return this.$store.state.records.filter(r => r.zone === this.zone)
    },
    filtered () {
      return this.filter === '' ? this.records : this.records.filter(r => r.domain.toLowerCase().includes(this.filter.toLowerCase()))
    }
  }
}
</script>