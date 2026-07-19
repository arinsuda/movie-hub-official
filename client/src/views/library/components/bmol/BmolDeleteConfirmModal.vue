<template>
  <ConfirmModal
    :model-value="modelValue"
    list-type="bmol_delete"
    :item-name="itemName"
    :confirm-disabled="confirmDisabled"
    @confirm="emit('confirm')"
    @cancel="emit('cancel')"
    @update:modelValue="emit('update:modelValue', $event)"
  />
</template>

<script setup lang="ts">
import ConfirmModal from "@/components/profile/components/ConfirmModal.vue"

defineProps<{
  modelValue: boolean
  itemName: string
  confirmDisabled: boolean
}>()

const emit = defineEmits<{
  (e: "confirm"): void
  (e: "cancel"): void
  (e: "update:modelValue", value: boolean): void
}>()
</script>
