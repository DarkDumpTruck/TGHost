<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { type Room, listRoom } from '@/api/actions'
import { ElButton } from 'element-plus';
import CreateRoomDialog from './CreateRoomDialog.vue';

const rooms = ref<Room[]>([])
const showCreateRoomDialog = ref(false)

const refetch = async () => {
  rooms.value = await listRoom()
}

onMounted(async () => {
  await refetch()
})
</script>

<template>
  <div class="space-y-4">
    <h1 class="text-center text-3xl font-bold underline">
      房间列表
    </h1>
    <div class="flex justify-end">
      <ElButton type="primary" plain @click="showCreateRoomDialog=true">
        创建房间
      </ElButton>
    </div>
    <CreateRoomDialog v-model="showCreateRoomDialog" @onConfirm="refetch" />
    <div v-for="room in rooms" :key="room.id">
      <div class="min-h-16 rounded-lg border-solid border-2 border-neutral-600">
        <div class="whitespace-pre-wrap text-lg m-4 space-x-2">
          <span class="font-bold text-lg/loose underline">{{ room.name }}</span>
          <p class="text-lg/loose">当前游戏：{{ room.gameName }} </p>
            <p class="text-lg/loose" v-for="player in room.players" :key="player.id">
              <router-link :to="`/game/${room.id}/${player.id}`">
                  {{ player.name }} 入口
              </router-link>
            </p>
        </div>
      </div>
    </div>
  </div>
</template>
