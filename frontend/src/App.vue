<template>
  <p v-if="connected">WS Connected sucessfully</p>
  <p v-else>No WS Connection</p>
  <div class="container mt-5">
        <div class="row">
            <div class="col-md-6 bg-light border">
                <h4 class="text-center">I/O Status</h4>
                <IoGrid :cards="inventory"></IoGrid>
            </div>
            <div class="col-md-6 bg-light border">
                <h4 class="text-center">Card Inventory</h4>
                <Inventory :data="inventory"></Inventory>
                <Monitor></Monitor>
                <h4 class="text-center">Log</h4>
                <Log></Log>
            </div>
        </div>
    </div>
  <div>
  </div>
</template>

<script>
import { mapState, mapWritableState } from "pinia";
import { SessionDataStore } from "./stores/session";

import Inventory from "./Inventory.vue";
import IoGrid from "./iogrid.vue";
import Log from "./Log.vue";
import Monitor from "./Monitor.vue"

export default {
  name: "App",
  setup() {
    const session = SessionDataStore();
    return { session };
  },
  components: {
    Inventory,
    IoGrid,
    Log,
    Monitor,
  },
  watch: {
    connected() {
      this.reconnect()
    },
  },
  mounted() {
    this.session.connect();
  },
  computed: {
    ...mapState(SessionDataStore, ["connected", "bytedata", "inventory"]),
    wsdata() {
      return JSON.stringify(this.bytedata)
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
