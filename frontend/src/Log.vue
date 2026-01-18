<template>
  <div>
    <label>Max Log Address <input v-model="maxlogaddress" type="number" step="1"/></label>
    <table class="" style="width: 100%">
      <thead>
        <tr>
          <th>Time</th>
          <th>Address</th>
          <th>Value</th>
        </tr>
      </thead>
     <tbody>
        <tr v-for="item in log">
          <td>{{ item.ts.toLocaleTimeString() }}</td>
          <td>{{ item.data.addr }}</td>
          <td>{{ item.data.value }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { SessionDataStore } from './stores/session';
import { mapState, mapWritableState } from 'pinia';

export default {
  name: "Log",
  setup() {
    const session = SessionDataStore();
    return { session };
  },
  computed: {
    ...mapState(SessionDataStore, ["log"]),
    ...mapWritableState(SessionDataStore, ["maxlogaddress"])
  },
};
</script>

<style>
</style>
