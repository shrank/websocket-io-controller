<template>
  <div class="row" id="MainApp" >
    <p v-if="connected">WS Connected sucessfully</p>
    <p v-else>No WS Connection</p>
    <textarea>
    {{ wsdata }}
    </textarea>
  </div>
</template>

<script>
import { mapState } from "pinia";
import { SessionDataStore } from "./stores/session";

export default {
  name: "App",
  setup() {
    const session = SessionDataStore();
    return { session };
  },
  components: {},
  data() {
    return {
    };
  },
  watch: {
    connected() {
    },
  },
  mounted() {
    this.session.connect();
  },
  computed: {
    ...mapState(SessionDataStore, ["connected", "data"]),
    wsdata() {
      return JSON.stringify(this.data)
    }
  },
  methods: {
    reconnect() {
      if(this.connected == true) {
        return
      }
      console.log("trying to reconnect")
      this.session.connect();
      setTimeout(this.reconnect, 2000)
    }
  },
};
</script>

<style>
</style>
