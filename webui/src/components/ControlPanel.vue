<template>
<div>
    <transition name="alert-fade">
      <div v-if="showAlert" class="container position-fixed alert-container">
        <div class="alert" :class="{ 'alert-success': status === 'ok', 'alert-danger': status === 'error' }" role="alert">
          {{ alertMessage }}
        </div>
      </div>
    </transition>
  <div class="container mt-3">
    <div class="row justify-content-between align-items-center"> 
      <div class="col-auto"><h1 class="site-header"><img src="@/assets/koala.png" alt="koala." style="height: 1.5em"/></h1></div>
      <div class="col-auto"><button class="btn btn-primary" @click="push" :disabled="applying">Apply changes</button></div>
    </div>
    <hr />
      <transition-group name="record-list" tag="div">
      <div v-for="rec in records" :key="records.indexOf(rec)" class="dns-record mb-2">
        <div class="row align-items-center m-2 p-md-1">
        <div class="col col-md-2 order-1 mb-3 mb-md-0">
          <span class="dns-record-type" @click="swap(rec)" :class="['dns-record-type-' + rec.type]">{{ rec.type }}</span>
        </div>
        <div class="col-md order-3 mb-1 mb-md-0">
          <input type="text" class="form-control" v-model="rec.domain">
        </div>
        <div class="col-md order-4 mb-1 mb-md-0">
          <input type="text" class="form-control" v-model="rec.data">
        </div>
        <div class="col-auto order-2 order-md-12 mb-3 mb-md-0">
          <button class="btn btn-outline-danger d-block" @click="del(rec)">Delete</button>
        </div>
        </div>
      </div>
      </transition-group>
      <div class="m-3">
      <div class="dns-record-add row p-3" @click="add">
        <div class="col-12 text-center">
          Add new record!
        </div>
      </div>
      </div>
  </div>
  </div>
</template>

<script>
import axios from "axios";

const alertTimeout = 3000; // 3 seconds timeout should be enough

export default {
  name: "ControlPanel",
  data() {
    return {
      records: [
        {
          type: "A",
          domain: "chatd",
          data: "192.168.10.130"
        }
      ],
      applying: false,
      showAlert: false,
      alertMessage: "",
      status: "success"
    };
  },
  created() {
    this.fetch();
  },
  methods: {
    fetch() {
      axios.get("/api/list").then(resp => {
        this.records = resp.data;
      });
    },
    push() {
      this.applying = true;
      axios.post("/api/apply", this.records).then(
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
        domain: "random-name",
        data: "192.168.1.1"
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
.site-header {
  font-weight: bold;
  font-size: 2em;
}
.dns-record {
  border: 1px solid #e9e9e9;
  transition: all 0.2s ease-in-out;
}
.dns-record:hover {
  transform: scale(1.015) translateY(-5px);
  box-shadow: 0px 3px 5px #ddd;
  background-color: #fcfcfc;
}
.dns-record-type {
  background-color: #2c3e50;
  color: #e9e9e9;
  font-weight: bold;
  display: block;
  padding: 0.5em;
  border-radius: 1em;
  min-width: 3em;
  text-align: center;
  cursor: pointer;
}
.dns-record-type-A {
  background-color: #1e88e5; /* blue 600 */
}
.dns-record-type-A:active {
  background-color: #0d47a1; /* blue 900 */
}
.dns-record-type-AAAA {
  background-color: #6746c3;
}
.dns-record-type-CNAME {
  background-color: #8e24aa; /* purple 600 */
}
.dns-record-type-CNAME:active {
  background-color: #4a148c; /* purple 900 */
}
.dns-record-add {
  color: #b2b2b2;
  border: dashed 2px #b2b2b2;
  cursor: pointer;
}
.record-list-enter-active,
.record-list-leave-active {
  transition: all 0.3s;
}
.record-list-enter {
  opacity: 0;
  transform: scaleY(0);
}
.record-list-leave-to {
  opacity: 0;
  transform: scaleY(0);
}
.alert-container {
  left: 50%;
  transform: translateX(-50%);
  z-index: 100;
}
.alert-fade-enter-active,
.alert-fade-leave-active {
  transition: opacity 0.5s;
}
.alert-fade-enter,
.alert-fade-leave-to {
  opacity: 0;
}
</style>
