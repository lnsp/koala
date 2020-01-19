<template>
  <div class="min-h-screen bg-gray-900 p-4">
    <transition name="alert-fade">
      <div v-if="showAlert" class="fixed w-screen mt-4 text-center">
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
    <div class="max-w-5xl mx-auto mb-6">
      <div class="py-3 flex flex-row justify-between items-center">
        <div class="flex-auto">
          <brand class="text-gray-100" />
        </div>
        <button
          class="bg-none font-medium text-gray-500 border border-gray-800 hover:border-indigo-700 hover:bg-indigo-700 hover:text-white focus:outline-none focus:shadow-outline rounded py-2 px-4"
          @click="push"
          :disabled="applying"
        >Apply changes</button>
      </div>
      <div class="bg-white shadow-md rounded mt-4 p-6">
        <div class="flex flex-row mb-6">
          <button class="text-white px-4 py-2 mx-3 rounded bg-indigo-600 hover:bg-indigo-700 shadow" @click="add">Add Record</button>
          <textInput class="flex-auto mx-3" v-model="filter" placeholder="Search for record" />
        </div>
        <recordList :records="records" :filter="filter" />
      </div>
    </div>
  </div>
</template>

<script>
import Brand from "../components/Brand";
import RecordList from "../components/RecordList";
import TextInput from '../components/TextInput';
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
      filter: "",
      axios: null,
    };
  },
  components: {
    brand: Brand,
    recordList: RecordList,
    textInput: TextInput,
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
      this.records.unshift({
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
