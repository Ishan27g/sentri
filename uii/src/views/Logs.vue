<script setup lang="ts">
import { ref } from "vue";

let fetchingContainers = ref(false);
let fetchingLogs = ref(false);
type container = { name: string, id: string }
type log = { prefix: string, color: string, text: string }

let serverPort = 4999;
let httpServer = "http://localhost:";
let wsServer = "ws://localhost:";
let logWsEndpoint = "/logs";
let DisplayLog = ref<log[]>([]);

let logConn: WebSocket = new WebSocket(wsServer + serverPort + logWsEndpoint);
logConn.onopen = () => {
  logConn.send(JSON.stringify("ok"));
};

logConn.onmessage = (event) => {
  DisplayLog.value = DisplayLog.value.concat(JSON.parse(event.data));
};

</script>

<template>
  <main>
    <div class="grid">
      <div class="col-11">
        <ScrollPanel style="width: 100%; height: 800px" class="scroll">
          <div class="grid grid-nogutter align-content-start" id="log" v-for="log in DisplayLog" :key="log.prefix">
            <div class="col-1" id="names" :style="{ color: log.color }">
              {{ log.prefix }}
            </div>
            <div class="col-11">
              {{ log.text }}
            </div>
          </div>
          <div>
          </div>
        </ScrollPanel>
      </div>
    </div>

  </main>
</template>

<style>
#log {
  font-family: Roboto, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-size: 15px;
  font-weight: 500;
  margin: auto;
  padding: auto;
}

#names {
  font-size: 16px;
}

#scroll {
  font-weight: 500;
}
</style>