import { computed, ref } from "vue";
import { i18n } from "./index";
import { SUPPORTED_LOCALES, type SupportedLocale } from "./types";
import { setLocale as apiSetLocale } from "./locale";

const isSwitchingLocale = ref(false);

export function useLocale() {
  const locale = computed({
    get() {
      return (i18n.global.locale as unknown as { value: SupportedLocale }).value;
    },
    set(val: SupportedLocale) {
      changeLocale(val);
    },
  });

  const isThai = computed(() => locale.value === "th");
  const isEnglish = computed(() => locale.value === "en");

  async function changeLocale(newLocale: SupportedLocale) {
    if (isSwitchingLocale.value) return false;
    isSwitchingLocale.value = true;
    try {
      return await apiSetLocale(newLocale);
    } finally {
      isSwitchingLocale.value = false;
    }
  }

  async function toggleLocale() {
    const nextLocale: SupportedLocale = locale.value === "th" ? "en" : "th";
    return await changeLocale(nextLocale);
  }

  return {
    locale,
    isThai,
    isEnglish,
    isSwitchingLocale,
    setLocale: changeLocale,
    toggleLocale,
    supportedLocales: SUPPORTED_LOCALES,
  };
}
