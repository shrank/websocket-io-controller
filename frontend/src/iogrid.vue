<template>
  <div>
    <table class="" style="width: auto">
      <thead>
        <tr>
          <th>Addr</th>
          <th>Mode</th>
          <th>Value</th>
          <th>Action</th>
          <th>Card</th>
          <th>Channel</th>
        </tr>
      </thead>
     <tbody>
        <tr v-for="item in allAddr" :key="item.addr">
          <td>{{ item.addr }}</td>
          <td>{{ item.mode }}</td>
          <td v-if="item.value > 0" style="background-color: lightgreen !important;">{{ item.value }}</td>
          <td v-else >{{ item.value }}</td>
          <td>
            <button v-if="['IN', 'OUT'].includes(item.mode)" @click="toggle(item)">Toggle</button>
            <button @click="monitor(item)">Monitor</button>
            </td>
          <td v-if="item.span > 0 " :class="item.status" :rowspan="item.span"><div class="text-topdown">Card #{{ item.slot }}</div></td>
          <td class="grid">{{ item.channel }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { SessionDataStore } from './stores/session';
import { mapState } from 'pinia';

export default {
  name: "IoGrid",
  setup() {
    const session = SessionDataStore();
    return { session };
  },
  props:{
    cards: Array
  },
  computed: {
    ...mapState(SessionDataStore, ["bytedata"]),
    max() {
      let m = 0
      for (let c of this.cards) {
        if(m < c.StartAddr + c.AddrCount){
          m = c.StartAddr + c.AddrCount
        }
      }
      return m
    },
    allAddr() {
      const res = []
      let i=0
      while(i < this.max){
        let card = null
        for (let c of this.cards) {
          if(c.StartAddr == i) {
            card = c
            break
          }
        }
        if(card == null) {
          res.push({
            addr: i,
            value: this.getValue(i),
            slot: "",
            span: 0,
            mode: "",
            status:"missing"
          })
          i++
        } else {
          res.push({
            addr: i,
            value: this.getValue(i),
            slot: card.BusAddr,
            span: card.AddrCount,
            mode: card.Mode,
            status: card.Status,
            channel: i - card.StartAddr
          })
          i++
          while( i < card.StartAddr + card.AddrCount) { 
            res.push({
              addr: i,
              value: this.getValue(i),
              slot: card.BusAddr,
              span: 0,
              mode: card.Mode,
              status: card.Status,
              channel: i - card.StartAddr
            })
            i++
          }
        }
      }
      return res
    }
  },
  methods: {
    toggle(item) {
      const update = {}
      update[item.addr] = (item.value + 1 ) % 2
      this.session.IoUpdate(update)
    },
    getValue(i) {
      return this.bytedata[i]
    },
    monitor(item) {
      this.session.setMonitor(item.addr)
    }
  }
};
</script>

<style>
.READY {
  background-color: lightgreen !important;
  border-color: darkgreen;
  border-style: solid;
  border-width: 1px;
}
.missing {
  background-color: gray !important;
}
.text-topdown {
  writing-mode: vertical-rl
}
.grid {
  border-color: #006400;
  border-style: solid;
  border-width: 1px;
}
</style>
