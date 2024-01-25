<script lang="ts" setup>
import { ref } from 'vue'
import { createRoom } from '@/api/actions'
const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits(['update:modelValue', 'onConfirm'])

const gameNames: Record<string,string> = {
  '偷袭': 'sneak-attack',
  '自定义游戏': 'custom',
}

const name = ref("")
const game = ref("")
const code = ref("")
const playerNum = ref(2)
const hidden = ref(false)

async function onGameChange (game: string) {
  name.value = game
  code.value = await fetch(`/assets/${gameNames[game]}.template`).then(resp => resp.text())
}

async function onConfirm () {
    await createRoom({name: name.value, code: code.value, playerNum: playerNum.value, hidden: hidden.value})
    emit('onConfirm')
    emit('update:modelValue', false)
}
</script>

<template>
    <el-dialog
      :model-value="modelValue"
      title="创建房间"
      align-center
      width="90%"
      @update:modelValue="(value: boolean) => $emit('update:modelValue', value)"
    >
    <div class="mx-4">
      <el-form label-position="top">
        <el-form-item label="游戏">
          <el-select v-model="game" @change="onGameChange" >
            <el-option v-for="k, g in gameNames" :key="k" :label="g" :value="g" />
          </el-select>
        </el-form-item>
        <el-form-item label="房间名">
        <el-input v-model="name" placeholder="房间名"></el-input>
        </el-form-item>
        <el-form-item label="游戏代码">
        <el-input v-model="code" placeholder="游戏代码" type="textarea" :autosize="{ minRows: 5, maxRows: 15 }" />
        </el-form-item>
        <el-form-item label="人数">
        <el-slider v-model="playerNum" show-input :min="2" :max="16" :step="1" />
        </el-form-item>
        <el-form-item label="隐藏">
        <el-switch v-model="hidden"></el-switch>
        </el-form-item>
      </el-form></div>
      <template #footer>
        <span>
          <el-button @click="$emit('update:modelValue', false)">取消</el-button>
          <el-button type="primary" @click="onConfirm">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
</template>