<template>
  <div class="container mt-3">
    <div class="row justify-content-between"> 
      <div class="col-auto"><h1 class="site-header">koala.</h1></div>
      <div class="col-auto"><button class="btn btn-primary" @click="push" :disabled="applying">Apply changes</button></div>
    </div>
    <hr />
    <div>
      <div v-for="rec in records" :key="rec.domain" class="row dns-record p-3 align-items-center mb-2">
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
      <div class="p-3 dns-record-add row mb-3" @click="add">
        <div class="col-12 text-center">
          Add new record!
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

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
      axios.post('/api/apply', this.records).then((resp) => {
        this.applying = false
      }, (error) => {
        console.log('apply request failed', error)
      })
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
  transition: transform 0.15s ease-in-out;
  transition-property: transform box-shadow background-color;
}
.dns-record-type:hover {
  transform: translateY(-5px);
  box-shadow: 0px 3px 5px #aaa;
}
.dns-record-type:active {
  transform: translateY(-3px);
  box-shadow: 0px 3px 5px #aaa, 0px 5px rgba(0,0,0,0.4) inset;
}
.dns-record-type-A {
  background-color: #311b92;
}
.dns-record-type-AAAA {
  background-color: #6746c3;
}
.dns-record-type-CNAME {
  background-color: #880e4f;
}
.dns-record-add {
  color: #b2b2b2;
  border: dashed 2px #b2b2b2;
  cursor: pointer;
}
</style>
