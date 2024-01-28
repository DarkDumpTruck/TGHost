<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { ElMessage, type ElSlider } from 'element-plus'
import { useRouteParams } from '@vueuse/router'
import { type PlayerStatus, getPlayerStatus, postPlayerInput } from '@/api/actions'

const roomID = useRouteParams<number>("roomId", 0, { transform: Number })
const playerID = useRouteParams<number>("playerId", 0, { transform: Number })

const playerStatus = ref<PlayerStatus>()
const inputDDL = ref<number>(0)
const inputValue = ref<string>("")
const sliderValue = ref<number>(0)
const sliderMin = ref<number>(0)
const sliderMax = ref<number>(0)
const sliderStep = ref<number>(0)

const playerStatusContainer = ref<HTMLDivElement>()

const websocketStatus = ref<string>("连接中")
const websocketProtocol = window.location.protocol === "https:" ? "wss" : "ws"
const websocketConn = ref<WebSocket>()

async function refetch() {
  const data = await getPlayerStatus({ roomId: roomID.value, playerId: playerID.value })
  playerStatus.value = data
  inputDDL.value = data.inputDDL
  playerStatusContainer.value?.scrollTo(0, playerStatusContainer.value.scrollHeight)
  if (data.inputType.startsWith("slider:")) {
    let range = data.inputType.split(":").map(Number)
    if (range.length > 2) {
      sliderMin.value = range[1]
      sliderMax.value = range[2]
      sliderValue.value = range[1]
    } else {
      sliderMin.value = 0
      sliderMax.value = range[1]
      sliderValue.value = 0
    }
    if (range.length > 3) {
      sliderStep.value = range[3]
    } else {
      sliderStep.value = 1
    }
  }
}

async function submit() {
  if (!playerStatus.value) {
    return
  }
  let msg = ""
  if (playerStatus.value.inputType.startsWith("slider")) {
    msg = String(sliderValue.value)
  } else if (playerStatus.value.inputType.startsWith("input")) {
    msg = inputValue.value
  }
  let resp = await postPlayerInput(
    {
      roomId: roomID.value,
      playerId: playerID.value,
      inputId: playerStatus.value.inputId,
      msg,
    })
  if (resp.status === "ok") {
    inputValue.value = ""
  }
  refetch()
}

function startWS() {
  let conn = new WebSocket(`${websocketProtocol}://${window.location.host}/ws/${roomID.value}/${playerID.value}`)
  conn.onopen = () => {
    websocketStatus.value = "已连接"
  }
  conn.onclose = () => {
    websocketStatus.value = "已关闭"
  }
  conn.onmessage = (event) => {
    const data = JSON.parse(event.data)
    if (data === "update") {
      refetch()
    }
    if (data.startsWith("alert:")) {
      let type = data.split(":")[1]
      let message = data.split(":")[2]
      ElMessage({ type, message })
    }
  }

  websocketConn.value = conn
}

function closeWS() {
  websocketConn.value?.close()
}

onMounted(() => {
  refetch()
  setInterval(() => {
    if (inputDDL.value > 0)
      inputDDL.value--
  }, 1000)
  startWS()
})

onUnmounted(() => {
  closeWS()
})

</script>

<template>
  <div class="space-y-2">
    <div class="flex flex-wrap space-x-4">
      <div class="flex-grow">
        <h1 class="text-3xl/loose font-bold underline select-none">
          {{ playerStatus?.gameName }}
        </h1>
      </div>
      <div>
        <p class="text-lg select-none">
          服务器状态：{{ websocketStatus }}
        </p>
      </div>
      <router-link to="/rooms" class="text-lg">
        返回房间列表
      </router-link>
    </div>
    <div class="min-h-64 rounded-lg border-solid border-2 border-neutral-600">
      <div ref="playerStatusContainer" class="whitespace-pre-wrap text-lg m-2 max-h-[640px] overflow-y-scroll" v-text="playerStatus?.gameStatus">
      </div>
    </div>
    <div class="space-y-2">
      <p class="text-lg">
        <span v-if="playerStatus?.gameRunning === false">
          游戏已结束。
        </span>
        <span v-else>
          {{ playerStatus?.inputMsg }}
          <span v-if="playerStatus?.inputDone === false">
            ，剩余 {{ inputDDL }} 秒
          </span>
          <span v-else>
            ，正在等待其他玩家。
          </span>
        </span>
      </p>
      <div v-show="playerStatus?.inputType?.startsWith('input')">
        <el-input v-model="inputValue" :disabled="playerStatus?.inputDone" @keyup.enter="submit" />
      </div>
      <div v-show="playerStatus?.inputType?.startsWith('slider')" class="mx-4">
        <el-slider v-model="sliderValue" :min="sliderMin" :max="sliderMax" :step="sliderStep"
          :disabled="playerStatus?.inputDone" show-input />
      </div>
      <el-button type="primary" @click="submit">
        提交
      </el-button>
    </div>
  </div>
</template>
