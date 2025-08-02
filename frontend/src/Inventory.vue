<template>
  <div>
    <table class="table">
      <thead>
        <tr>
          <th>Type</th>
          <th>Start Address</th>
          <th>Address Count</th>
          <th>Word Size</th>
          <th>Bus Addr</th>
          <th>Mode</th>
        </tr>
      </thead>
     <tbody>
        <tr v-for="item in data" :key="item.StartAddr">
          <td>{{ item.Type }}</td>
          <td>{{ item.StartAddr }}</td>
          <td>{{ item.AddrCount }}</td>
          <td>{{ item.WordSize }}</td>
          <td>{{ item.BusAddr }}</td>
          <td>{{ item.Mode }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { ssrImportMetaKey } from 'vite/module-runner';

export default {
  name: "Inventory",
  props:{
    data: Array
  },
};
</script>

<style>
</style>
