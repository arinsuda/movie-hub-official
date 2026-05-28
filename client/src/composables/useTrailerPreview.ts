// client/src/composables/useTrailerPreview.ts
import { ref, computed } from "vue";
import type { Video } from "@/types";

// ─── Types ────────────────────────────────────────────────────────────────────

export type TrailerPriority = "trailer" | "teaser" | "clip" | "other";

export interface ResolvedTrailer {
  key: string;
  name: string;
  type: string;
  priority: TrailerPriority;
  embedUrl: string;
}

export interface PopupPosition {
  top: number;
  left: number;
  width: number;
  transformOrigin: "top left" | "top center" | "top right";
}

// ─── resolveTrailer ───────────────────────────────────────────────────────────

const TYPE_SCORE: Record<string, number> = {
  Trailer: 0,
  Teaser: 1,
  Clip: 2,
  Featurette: 3,
  "Behind the Scenes": 4,
  Bloopers: 5,
};

const PRIORITY_MAP: Record<string, TrailerPriority> = {
  Trailer: "trailer",
  Teaser: "teaser",
  Clip: "clip",
};

function buildEmbedUrl(key: string): string {
  const p = new URLSearchParams({
    autoplay: "1",
    mute: "1",
    controls: "0",
    playsinline: "1",
    rel: "0",
    loop: "1",
    playlist: key, // required for loop to function
    modestbranding: "1",
    iv_load_policy: "3", // hide video annotations
    disablekb: "1",
    fs: "0",
  });
  return `https://www.youtube.com/embed/${key}?${p.toString()}`;
}

export function resolveTrailer(
  videos: Video[] | undefined | null,
): ResolvedTrailer | null {
  if (!videos?.length) return null;

  const yt = videos.filter((v) => v.site === "YouTube");
  if (!yt.length) return null;

  const sorted = [...yt].sort((a, b) => {
    if (a.official !== b.official) return a.official ? -1 : 1;
    const aS = TYPE_SCORE[a.type] ?? 99;
    const bS = TYPE_SCORE[b.type] ?? 99;
    return aS - bS;
  });

  const best = sorted[0];
  if (!best) return null; // ✅ guard สำหรับ TypeScript

  return {
    key: best.key,
    name: best.name,
    type: best.type,
    priority: PRIORITY_MAP[best.type] ?? "other",
    embedUrl: buildEmbedUrl(best.key),
  };
}

// ─── useTrailerPreview ────────────────────────────────────────────────────────

interface Options {
  /** How long after hover before iframe is injected into DOM (default: 500ms) */
  mountDelay?: number;
}

export function useTrailerPreview(options: Options = {}) {
  const { mountDelay = 500 } = options;

  const isIframeMounted = ref(false);
  const isIframeLoaded = ref(false);

  let timer: ReturnType<typeof setTimeout> | null = null;

  function scheduleMount() {
    cancelMount();
    timer = setTimeout(() => {
      isIframeMounted.value = true;
    }, mountDelay);
  }

  function cancelMount() {
    if (timer !== null) {
      clearTimeout(timer);
      timer = null;
    }
  }

  function unmount() {
    cancelMount();
    isIframeMounted.value = false;
    isIframeLoaded.value = false;
  }

  function onIframeLoad() {
    isIframeLoaded.value = true;
  }

  // Derived state for template convenience
  const showSkeleton = computed(
    () => isIframeMounted.value && !isIframeLoaded.value,
  );
  const showFallback = computed(() => !isIframeMounted.value);

  return {
    isIframeMounted,
    isIframeLoaded,
    showSkeleton,
    showFallback,
    scheduleMount,
    cancelMount,
    unmount,
    onIframeLoad,
  };
}

// ─── usePopupPosition ─────────────────────────────────────────────────────────

const POPUP_WIDTH = 300; // px — match .hover-popup width in HomeView
const POPUP_OFFSET_Y = -8; // shift popup up slightly relative to card top

export function usePopupPosition() {
  const position = ref<PopupPosition | null>(null);

  function calculate(cardEl: HTMLElement): void {
    const rect = cardEl.getBoundingClientRect();
    const vw = window.innerWidth;
    const scrollY = window.scrollY;

    const cardCenterX = rect.left + rect.width / 2;
    let left = cardCenterX - POPUP_WIDTH / 2;

    type Side = "left" | "center" | "right";
    let side: Side = "center";

    if (left < 12) {
      left = 12;
      side = "left";
    } else if (left + POPUP_WIDTH > vw - 12) {
      left = vw - POPUP_WIDTH - 12;
      side = "right";
    }

    const originMap: Record<Side, PopupPosition["transformOrigin"]> = {
      left: "top left",
      center: "top center",
      right: "top right",
    };

    position.value = {
      top: rect.top + scrollY + POPUP_OFFSET_Y,
      left,
      width: POPUP_WIDTH,
      transformOrigin: originMap[side],
    };
  }

  function clear() {
    position.value = null;
  }

  return { position, calculate, clear };
}
