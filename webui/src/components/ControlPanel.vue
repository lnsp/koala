<template>
  <div class="min-h-screen bg-gray-900 p-4">
    <transition name="alert-fade">
      <div v-if="showAlert"
           class="z-50 fixed w-full -ml-4 mt-4 text-center">
        <div class="items-center p-2 rounded-full inline-flex"
             :class="{ 'text-green-100': status === 'ok', 'bg-green-700': status === 'ok', 'text-red-100': status === 'error', 'bg-red-700': status === 'error'}">
          <div class="flex rounded-full px-3 py-1 mr-4 uppercase font-bold text-xs"
               :class="{ 'bg-green-500': status === 'ok', 'bg-red-500': status === 'error'}">{{ status }}</div>
          <div class="mr-3 text-xs sm:text-base font-semibold flex-auto">{{ alertMessage }}</div>
        </div>
      </div>
    </transition>
    <div class="max-w-5xl mx-auto mb-6">
      <div class="py-3 flex flex-row justify-between items-center">
        <div class="flex-auto">
          <brand-header class="text-gray-100" />
        </div>
        <button class="bg-none flex items-center font-medium text-gray-500 border border-gray-800 hover:border-indigo-700 hover:bg-indigo-700 hover:text-white focus:outline-none focus:shadow-outline rounded py-2 px-4"
                @click="push"
                :disabled="applying">
          <svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-white"
               xmlns="http://www.w3.org/2000/svg"
               fill="none"
               viewBox="0 0 24 24"
               v-if="applying">
            <circle class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"></circle>
            <path class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Apply changes</button>
      </div>
      <div class="bg-white shadow-md rounded mt-4">
        <TabGroup>
          <TabList class="flex gap-4 sm:mx-6 mx-3 pt-3 border-b border-gray-400">
            <Tab as="template"
                 v-for="zone in zones"
                 :key="zone"
                 v-slot="{ selected }">
              <button :class="selected ? ['text-indigo-600', 'border-indigo-500'] : ['text-gray-400', 'border-transparent']"
                      class="font-medium px-3 py-2 border-b-2">{{ zone }}</button>
            </Tab>
          </TabList>
          <TabPanels>
            <TabPanel v-for="zone in zones"
                      :key="zone">
              <record-list :zone="zone" />
            </TabPanel>
          </TabPanels>
        </TabGroup>
      </div>
    </div>
  </div>
</template>

<script>
import BrandHeader from "../components/BrandHeader";
import RecordList from "../components/RecordList";
import { TabGroup, TabList, Tab, TabPanels, TabPanel } from "@headlessui/vue";
import { inject, ref, computed, onMounted } from "vue";
import { useStore } from "vuex";

const alertTimeout = 3000; // 3 seconds timeout should be enough

export default {
  components: {
    BrandHeader,
    RecordList,
    TabGroup,
    TabList,
    Tab,
    TabPanels,
    TabPanel,
  },
  setup() {
    const axios = inject("axios");
    const store = useStore();

    const applying = ref(false);
    const alertMessage = ref("");
    const status = ref("ok");
    const showAlert = ref(false);
    const activeZone = ref("");

    const showSuccess = (msg) => {
      status.value = "ok";
      alertMessage.value = msg;
      showAlert.value = true;
      setTimeout(() => (showAlert.value = false), alertTimeout);
    };

    const showError = (err) => {
      status.value = "error";
      alertMessage.value = err;
      showAlert.value = true;
      setTimeout(() => (showAlert.value = false), alertTimeout);
    };

    const updateZones = async () => {
      let response = await axios.get("/");

      // store list of zones
      store.commit("setZones", { zones: response.data });

      // set active zone
      activeZone.value = response.data[0];
    };

    const updateRecords = async (zone) => {
      let response = await axios.get("/list", { params: { zone: zone } });

      // store list of records for zone
      store.commit("setRecords", { zone: zone, records: response.data });
    };

    const fetch = async () => {
      try {
        await updateZones();

        for (let index in store.state.zones) {
          await updateRecords(store.state.zones[index]);
        }
      } catch (err) {
        showError("Sorry, something has gone wrong.");
        console.log(err);
      }
    };

    const pushZone = async (zone) => {
      const records = store.state.records.filter((r) => r.zone == zone);
      await axios.post("/apply", records, { params: { zone: zone } });
    };

    const push = async () => {
      applying.value = true;

      try {
        for (const zone of store.state.zones) {
          await pushZone(zone);
        }
        showSuccess("Your configuration change has been applied.");
      } catch (err) {
        showError("Sorry, something has gone wrong.");
      } finally {
        applying.value = false;
        fetch();
      }
    };

    onMounted(() => {
      fetch();
    });

    const zones = computed(() => store.state.zones);

    return {
      alertMessage,
      zones,
      status,
      push,
      applying,
      showAlert,
    };
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.record-list-enter-active,
.record-list-leave-active {
  transition: all 0.15s;
}
.record-list-enter {
  opacity: 0;
}
.record-list-leave-to {
  opacity: 0;
}
.alert-fade-enter-active,
.alert-fade-leave-active {
  transition: opacity 0.25s;
}
.alert-fade-enter,
.alert-fade-leave-to {
  opacity: 0;
}
</style>
