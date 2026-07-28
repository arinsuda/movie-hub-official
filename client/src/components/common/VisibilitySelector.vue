<script setup lang="ts">
import { ref, computed } from "vue";
import { useI18n } from "vue-i18n";
import { Globe, Users, Lock, ChevronDown } from "lucide-vue-next";
import { useClickOutside } from "@/composables/useClickOutside";
import type { Visibility } from "@/types";

const props = withDefaults(
  defineProps<{
    modelValue?: Visibility;
    disabled?: boolean;
    size?: "sm" | "md" | "lg";
  }>(),
  {
    modelValue: "public",
    disabled: false,
    size: "md",
  }
);

const emit = defineEmits<{
  (e: "update:modelValue", value: Visibility): void;
}>();

const { t } = useI18n();

const isOpen = ref(false);
const selectorRef = ref<HTMLElement | null>(null);

useClickOutside(selectorRef, () => {
  isOpen.value = false;
});

const options = computed(() => [
  {
    value: "public" as Visibility,
    icon: Globe,
    title: t("visibility.public"),
    desc: t("visibility.publicDesc"),
  },
  {
    value: "followers" as Visibility,
    icon: Users,
    title: t("visibility.followers"),
    desc: t("visibility.followersDesc"),
  },
  {
    value: "private" as Visibility,
    icon: Lock,
    title: t("visibility.private"),
    desc: t("visibility.privateDesc"),
  },
]);

const currentOption = computed(() => {
  const found = options.value.find((opt) => opt.value === props.modelValue);
  return found ?? options.value[0] ?? {
    value: "public" as Visibility,
    icon: Globe,
    title: t("visibility.public"),
    desc: t("visibility.publicDesc"),
  };
});

const toggleMenu = () => {
  if (props.disabled) return;
  isOpen.value = !isOpen.value;
};

const selectOption = (value: Visibility) => {
  if (props.disabled) return;
  emit("update:modelValue", value);
  isOpen.value = false;
};

const handleKeyDown = (e: KeyboardEvent) => {
  if (props.disabled) return;
  
  if (e.key === "Escape") {
    isOpen.value = false;
  } else if (e.key === "Enter" || e.key === " ") {
    e.preventDefault();
    toggleMenu();
  } else if (isOpen.value) {
    const currentIndex = options.value.findIndex(o => o.value === props.modelValue);
    if (e.key === "ArrowDown") {
      e.preventDefault();
      const nextIndex = (currentIndex + 1) % options.value.length;
      const targetOpt = options.value[nextIndex];
      if (targetOpt) {
        selectOption(targetOpt.value);
      }
      isOpen.value = true; // keep open while navigating
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      const prevIndex = (currentIndex - 1 + options.value.length) % options.value.length;
      const targetOpt = options.value[prevIndex];
      if (targetOpt) {
        selectOption(targetOpt.value);
      }
      isOpen.value = true;
    }
  }
};

const sizeClasses = computed(() => {
  switch (props.size) {
    case "sm":
      return "px-2 py-1 text-xs";
    case "lg":
      return "px-4 py-2.5 text-base";
    case "md":
    default:
      return "px-3 py-1.5 text-sm";
  }
});
</script>

<template>
  <div class="relative inline-block" ref="selectorRef">
    <button
      type="button"
      :class="[
        'flex items-center gap-2 rounded-md border border-neutral-700 bg-neutral-900 font-medium text-neutral-200 transition-colors hover:bg-neutral-800 hover:text-white focus:outline-none focus:ring-2 focus:ring-primary-500',
        sizeClasses,
        disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'
      ]"
      @click="toggleMenu"
      @keydown="handleKeyDown"
      :aria-expanded="isOpen"
      aria-haspopup="listbox"
      :disabled="disabled"
    >
      <component :is="currentOption.icon" class="h-4 w-4" />
      <span>{{ currentOption.title }}</span>
      <ChevronDown class="ml-1 h-4 w-4 opacity-50" />
    </button>

    <Transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="transform scale-95 opacity-0"
      enter-to-class="transform scale-100 opacity-100"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="transform scale-100 opacity-100"
      leave-to-class="transform scale-95 opacity-0"
    >
      <div
        v-if="isOpen"
        class="absolute right-0 z-50 mt-1 w-64 origin-top-right rounded-md border border-neutral-700 bg-neutral-900 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none"
        role="listbox"
        tabindex="-1"
      >
        <div class="py-1">
          <button
            v-for="option in options"
            :key="option.value"
            @click="selectOption(option.value)"
            type="button"
            class="flex w-full items-start gap-3 px-4 py-2 text-left hover:bg-neutral-800"
            :class="[
              modelValue === option.value ? 'bg-neutral-800/50 text-primary-400' : 'text-neutral-300'
            ]"
            role="option"
            :aria-selected="modelValue === option.value"
          >
            <component :is="option.icon" class="mt-0.5 h-4 w-4 shrink-0" />
            <div class="flex flex-col">
              <span class="text-sm font-medium">{{ option.title }}</span>
              <span class="text-xs text-neutral-400">{{ option.desc }}</span>
            </div>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>
