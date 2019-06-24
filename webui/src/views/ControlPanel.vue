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
      <div class="col-auto"><h1 class="site-header"><img src="@/assets/koala.png" alt="koala." style="height: 1em"/></h1></div>
      <div class="col-auto"><button class="btn btn-koala" @click="push" :disabled="applying">Apply changes</button></div>
    </div>
    <hr />
      <transition-group name="record-list" tag="div">
      <div v-for="rec in records" :key="records.indexOf(rec)" class="dns-record mb-2">
        <div class="row align-items-center m-2 p-md-1">
        <div class="col col-md-2 order-1 mb-3 mb-md-0">
          <span class="dns-record-type" @click="swap(rec)" :class="['dns-record-type-' + rec.type]">{{ rec.type }}</span>
        </div>
        <div class="col-md order-3 mb-1 mb-md-0">
          <input type="text" class="form-control" placeholder="Hostname" v-model="rec.domain">
        </div>
        <div class="col-md order-4 mb-1 mb-md-0">
          <input type="text" class="form-control" placeholder="Address" v-model="rec.data">
        </div>
        <div class="col-auto order-2 order-md-12 mb-3 mb-md-0">
          <button class="btn btn-danger d-block" @click="del(rec)">&#x1f5d9;</button>
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
  props: ['rootAPI'],
  data() {
    return {
      records: [],
      applying: false,
      showAlert: false,
      alertMessage: "",
      status: "success",
      axios: null,
    };
  },
  beforeCreate () {
    document.documentElement.className = 'controlPanel';
    document.body.className = 'controlPanel';
  },
  mounted() {
    this.axios = axios.create({
      baseURL: this.rootAPI,
      headers: {'Authorization': 'Bearer ' + localStorage.token},
    });
    this.fetch();
  },
  methods: {
    fetch() {
      this.axios.get('/list')
      .then(resp => {
        this.records = resp.data;
      })
      .catch(err => {
        if (err.response.status === 401) {
          this.$router.push('/auth');
        }
      });
    },
    push() {
      this.applying = true;
      this.axios.post('/apply', this.records)
      .then(
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
.site-header {
  font-weight: bold;
  font-size: 2em;
}
.dns-record {
  border: 2px solid #e9e9e9;
  transition: all 0.2s ease-in-out;

}
.dns-record:hover {
  border: 2px solid #4a138c;
  background-color: #fff;
}
.dns-record-type {
  background-color: #2c3e50;
  color: #e9e9e9;
  font-weight: bold;
  display: block;
  padding: 0.5em;
  border-radius: 0.2em;
  min-width: 3em;
  text-align: center;
  cursor: pointer;

  /*
  Disable text highlighting on click.
  SOURCE: https://stackoverflow.com/questions/826782/how-to-disable-text-selection-highlighting
  */
  -webkit-touch-callout: none; /* iOS Safari */
  -webkit-user-select: none; /* Safari */
  -khtml-user-select: none; /* Konqueror HTML */
  -moz-user-select: none; /* Firefox */
  -ms-user-select: none; /* Internet Explorer/Edge */
  user-select: none; /* Non-prefixed version, currently supported by Chrome and Opera */
}
input:focus {
  border: 2px solid #4a138c;
}
.dns-record-type-A {
  background-color: #4a148c; /* blue 600 */
}
.dns-record-type-A:active {
  background-color: #7b1fa2; /* blue 900 */
}
.dns-record-type-AAAA {
  background-color: #5e35b1;
}
.dns-record-type-CNAME {
  background-color: #e91e63; /* purple 900 */
}
.dns-record-type-CNAME:active {
  background-color: #c2185b; /* purple 600 */
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
}
.record-list-leave-to {
  opacity: 0;
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
