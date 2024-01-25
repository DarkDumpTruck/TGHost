<script lang="ts" setup>
import { useClipboard, useTimeout } from '@vueuse/core'

const props = defineProps<{ text: string }>()
const emit = defineEmits(['onCopy'])

const { copy } = useClipboard()
const { start, isPending } = useTimeout(1000, { controls: true, immediate: false })

function handleCopy () {
  start()
  copy(props.text)
  emit('onCopy')
}

</script>

<template>
  <div class="group relative flex pr-5">
    <div class="px-2">
      <slot />
    </div>
    <div
      class="absolute left-0 top-0 hidden cursor-pointer items-center rounded bg-gray-200 px-2 group-hover:flex space-x-1.5"
      @click="handleCopy()"
    >
      <div class="select-none">
        <slot />
      </div>
      <div :class="{ 'i-mdi:content-copy': !isPending, 'i-mdi:check text-green-600': isPending }" />
    </div>
  </div>
</template>
