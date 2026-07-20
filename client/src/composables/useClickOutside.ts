import { onBeforeUnmount, onMounted, type Ref } from "vue";

export function useClickOutside(
  target: Ref<HTMLElement | null>,
  handler: () => void,
) {
  function onClick(e: MouseEvent) {
    if (target.value && !target.value.contains(e.target as Node)) {
      handler();
    }
  }

  onMounted(() => document.addEventListener("click", onClick, true));
  onBeforeUnmount(() => document.removeEventListener("click", onClick, true));
}
