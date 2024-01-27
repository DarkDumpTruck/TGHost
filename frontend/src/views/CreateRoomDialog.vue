<script lang="ts" setup>
import { ref } from 'vue'
import { createRoom } from '@/api/actions'
const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits(['update:modelValue', 'onConfirm'])

const gameNames: Record<string,string> = {
  '偷袭': 'sneak-attack.js',
  '动物园捐赠': 'zoo-donation.js',
  '自定义游戏': 'custom.template',
}

const name = ref("")
const game = ref("")
const showCode = ref(false)
const code = ref("")
const totalNum = ref(2)
const playerNum = ref(2)
const botNum = ref(0)
const judgeNeeded = ref(false)
const hidden = ref(false)

async function onGameChange (game: string) {
  name.value = game
  code.value = await fetch(`/assets/${gameNames[game]}`).then(resp => resp.text())
  if(game == '自定义游戏') {
    showCode.value = true
  }
  onCodeChange(code.value)
}

function onCodeChange (code: string) {
  for(let line of code.split('\n')) {
    if(line.startsWith('//!player=')) {
      let num = parseInt(line.split('=')[1])
      if (Number.isInteger(num)) {
        totalNum.value = num
        playerNum.value = num
        botNum.value = 0
      }
    }
    if(line.startsWith('//!judge=')) {
      judgeNeeded.value = line[9] == '1'
    }
  }
}

function onPlayerNumChange (num: number) {
  if (num > totalNum.value) {
    playerNum.value = totalNum.value
  }
  botNum.value = totalNum.value - playerNum.value
}

async function onConfirm () {
    await createRoom({name: name.value, code: code.value, playerNum: playerNum.value, hidden: hidden.value, botNum: botNum.value, judgeNum: judgeNeeded.value?1:0})
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
        <el-checkbox v-model="showCode">显示游戏代码</el-checkbox>
        <el-form-item label="游戏代码" v-if="showCode">
        <el-input
          v-model="code"
          placeholder="游戏代码"
          type="textarea"
          :autosize="{ minRows: 5, maxRows: 15 }"
          @input="onCodeChange"/>
        </el-form-item>
        <el-form-item label="人数">
          <el-slider v-model="playerNum" show-input :min="1" :max="16" :step="1" @input="onPlayerNumChange"/>
        </el-form-item>
        <el-form-item label="bot数量">
          <el-slider v-model="botNum" show-input :min="0" :max="16" :step="1" />
        </el-form-item>
        <el-checkbox v-model="judgeNeeded">需要裁判</el-checkbox>
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