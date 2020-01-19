<template>
  <div class="h-screen w-screen bg-gray-200 p-4">
    <transition name="alert-fade">
      <div v-if="showAlert" class="fixed w-screen top-10 text-center">
        <div
          class="items-center p-2 rounded-full inline-flex"
          :class="{ 'text-green-100': status === 'ok', 'bg-green-700': status === 'ok', 'text-red-100': status === 'error', 'bg-red-700': status === 'error'}"
        >
          <div
            class="flex rounded-full px-3 py-1 mr-4 uppercase font-bold text-xs"
            :class="{ 'bg-green-500': status === 'ok', 'bg-red-500': status === 'error'}"
          >{{ status }}</div>
          <div class="mr-3 font-semibold flex-auto">{{ alertMessage }}</div>
        </div>
      </div>
    </transition>
    <div class="max-w-5xl mx-auto">
      <div class="pb-3 flex flex-row justify-between items-center border-b border-gray-300">
        <div class="flex-auto">
          <brand />
        </div>
        <button
          class="bg-none font-medium text-indigo-600 border border-indigo-600 hover:bg-indigo-700 hover:text-white focus:outline-none focus:shadow-outline rounded py-2 px-4"
          @click="push"
          :disabled="applying"
        >Apply changes</button>
      </div>
      <div class="bg-white shadow-md rounded mt-4 p-6">
        <recordList :records="records" />
        <div class="flex justify-center">
          <button class="text-white w-12 h-12 font-black rounded-full bg-indigo-600 shadow" @click="add">&#xFF0B;</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Brand from "../components/Brand";
import RecordList from "../components/RecordList";
import axios from "axios";

const alertTimeout = 3000; // 3 seconds timeout should be enough

export default {
  name: "ControlPanel",
  props: ["apiBaseURL"],
  data() {
    return {
      records: [],
      applying: false,
      showAlert: false,
      alertMessage: "",
      status: "success",
      axios: null
    };
  },
  components: {
    brand: Brand,
    recordList: RecordList
  },
  beforeCreate() {
    document.documentElement.className = "controlPanel";
    document.body.className = "controlPanel";
  },
  mounted() {
    this.axios = axios.create({
      baseURL: this.apiBaseURL,
      headers: { Authorization: "Bearer " + localStorage.token }
    });
    this.fetch();
  },
  methods: {
    fetch() {
      this.axios
        .get("/list")
        .then(resp => {
          this.records = resp.data;
        })
        .catch(err => {
          if (err.response.status === 401) {
            this.$router.push("/auth");
          }
        });
    },
    push() {
      this.applying = true;
      this.axios.post("/apply", this.records).then(
        () => {
          this.applying = false;
          this.showSuccess("Your configuration change has been applied.");
        },
        () => {
          this.applying = false;
          this.showError("Sorry, we could not contact the server.");
        }
      );
    },
    showSuccess(msg) {
      this.status = "ok";
      this.alertMessage = msg;
      this.showAlert = true;
      setTimeout(() => (this.showAlert = false), alertTimeout);
    },
    showError(err) {
      this.status = "error";
      this.alertMessage = err;
      this.showAlert = true;
      setTimeout(() => (this.showAlert = false), alertTimeout);
    },
    add() {
      this.records.push({
        type: "A",
        domain: "",
        data: ""
      });
    },
    swap(rec) {
      if (rec.type === "A") rec.type = "CNAME";
      else rec.type = "A";
    },
    del(rec) {
      var index = this.records.indexOf(rec);
      if (index > -1) {
        this.records.splice(index, 1);
      }
    }
  }
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
