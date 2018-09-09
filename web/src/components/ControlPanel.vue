<template>
  <div class="container mt-3">
    <transition name="alert-fade">
      <div class="alert" :class="{ 'alert-success': status === 'ok', 'alert-danger': status === 'error' }" role="alert" v-if="showAlert">
        {{ alertMessage }}
      </div>
    </transition>
    <div class="row justify-content-between align-items-center"> 
      <div class="col-auto"><h1 class="site-header"><img src="@/assets/koala.png" alt="koala." style="height: 1.5em"/></h1></div>
      <div class="col-auto"><button class="btn btn-primary" @click="push" :disabled="applying">Apply changes</button></div>
    </div>
    <hr />
      <transition-group name="record-list" tag="div">
      <div v-for="rec in records" :key="records.indexOf(rec)" class="row dns-record p-3 align-items-center mb-2">
        <div class="col-2">
          <span class="dns-record-type" @click="swap(rec)" :class="['dns-record-type-' + rec.type]">{{ rec.type }}</span>
        </div>
        <div class="col-4">
          <input type="text" class="form-control" v-model="rec.domain">
        </div>
        <div class="col-5">
          <input type="text" class="form-control" v-model="rec.data">
        </div>
        <div class="col-1">
          <button class="btn btn-outline-danger" @click="del(rec)">Delete</button>
        </div>
      </div>
      </transition-group>
      <div class="p-3 dns-record-add row mb-3" @click="add">
        <div class="col-12 text-center">
          Add new record!
        </div>
      </div>
  </div>
</template>

<script>
import axios from 'axios';

const alertTimeout = 3000 // 3 seconds timeout should be enough

export default {
  name: 'ControlPanel',
  data () {
    return {
      records: [{
        type: 'A',
        domain: 'chatd',
        data: '192.168.10.130'
      }],
      applying: false,
      showAlert: false,
      alertMessage: '',
      status: 'success',
    }
  },
  created () {
    this.fetch()
  },
  methods: {
    fetch () {
      axios.get('/api/list').then((resp) => {
        this.records = resp.data
      })
    },
    push () {
      this.applying = true
      axios.post('/api/apply', this.records).then(() => {
        this.applying = false
        this.showSuccess('Your configuration change has been applied.')
      }, (err) => {
        console.log('apply request failed', err)
        this.applying = false
        this.showError('Sorry, we could not contact the server.')
      })
    },
    showSuccess (msg) {
      this.status = 'ok'
      this.alertMessage = msg
      this.showAlert = true
      setTimeout(() => this.showAlert = false, alertTimeout)
    },
    showError (err) {
      this.status = 'error'
      this.alertMessage = err
      this.showAlert = true
      setTimeout(() => this.showAlert = false, alertTimeout)
    },
    add () {
      this.records.push({
        type: 'A',
        domain: 'random-name',
        data: '192.168.1.1',
      })
    },
    swap (rec) {
      if (rec.type === 'A') rec.type = 'CNAME';
      else rec.type = 'A';
    },
    del (rec) {
      var index = this.records.indexOf(rec)
      if (index > -1) {
        this.records.splice(index, 1)
      }
    },
  }
}
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
  display: inline-block;
  padding: 0.5em;
  border-radius: 1em;
  min-width: 3em;
  text-align: center;
  cursor: pointer;
}
.dns-record-type-A {
  background-color: #1E88E5; /* blue 600 */
}
.dns-record-type-A:active {
  background-color: #0D47A1; /* blue 900 */
}
.dns-record-type-AAAA {
  background-color: #6746c3;
}
.dns-record-type-CNAME {
  background-color: #8E24AA; /* purple 600 */
}
.dns-record-type-CNAME:active {
  background-color: #4A148C; /* purple 900 */
}
.dns-record-add {
  color: #b2b2b2;
  border: dashed 2px #b2b2b2;
  cursor: pointer;
}
.alert-fade-enter-active, .alert-fade-leave-active {
  transition: all 0.5s;
}
.alert-fade-enter, .alert-fade-leave-to {
  opacity: 0;
  transform: translateY(-50px);
}
.record-list-enter-active, .record-list-leave-active {
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
</style>
