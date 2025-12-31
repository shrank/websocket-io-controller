<template>
  <div>
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
import { mapState } from 'pinia';

export default {
  name: "Log",
  setup() {
    const session = SessionDataStore();
    return { session };
  },
  computed: {
    ...mapState(SessionDataStore, ["log"]),
  },
};
</script>

<style>
</style>
