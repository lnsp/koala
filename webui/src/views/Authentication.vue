<template>
  <div class="back-panel">
    <div class="auth-panel">
      <div class="site-header">
        <h1>
          <img src="@/assets/koala.png" alt="koala." style="height: 1em" />
        </h1>
        <div class="text-muted text-bold">DNS zone editor</div>
      </div>
      <hr />
      <form @submit.prevent="finish()">
        <div class="form-group">
          <input type="text" placeholder="Token" class="form-control" v-model="token" />
        </div>
        <button type="submit" class="btn btn-koala">Sign in</button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.back-panel {
  background: none;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}
.site-header {
  margin-left: 2em;
  margin-right: 2em;
}
.auth-panel {
  background-color: #fff;
  width: min-content;
  padding: 1em;
  border-radius: 5px;
}
</style>

<style>
html.authentication,
body.authentication {
  /*
  background: rgb(42,12,87);
  background: linear-gradient(220deg, rgba(42,12,87,1) 0%, rgba(117,72,232,1) 100%);
  */
  /*background: rgb(200, 200, 200);*/
  background-color: rgba(117, 72, 232, 1);
}
</style>

<script>
export default {
  name: "Authentication",
  data() {
    return {
      token: localStorage.token
    };
  },
  beforeCreate() {
    document.documentElement.className = "authentication";
    document.body.className = "authentication";
  },
  methods: {
    finish: function() {
      localStorage.token = this.token;
      this.$router.push("/");
    }
  }
};
</script>