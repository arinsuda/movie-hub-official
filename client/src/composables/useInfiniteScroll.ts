import { onBeforeUnmount, onMounted, watch, type Ref } from "vue";

interface UseInfiniteScrollOptions {
  disabled?: Ref<boolean>;
  rootMargin?: string;
}

// เรียก callback เมื่อ sentinel element เลื่อนเข้ามาในจอ ใช้ทำ infinite scroll ของฟีด
export function useInfiniteScroll(
  target: Ref<HTMLElement | null>,
  callback: () => void,
  options: UseInfiniteScrollOptions = {},
) {
  let observer: IntersectionObserver | null = null;

  function observe() {
    if (!target.value || observer) return;
    observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting && !options.disabled?.value) {
          callback();
        }
      },
      { rootMargin: options.rootMargin ?? "200px" },
    );
    observer.observe(target.value);
  }

  onMounted(observe);
  // เผื่อ sentinel ถูก mount ทีหลัง (เช่นอยู่หลัง v-else ที่ render ช้ากว่า)
  watch(target, (el) => {
    if (el && !observer) observe();
  });

  onBeforeUnmount(() => {
    observer?.disconnect();
    observer = null;
  });
}
