import { createI18n } from "vue-i18n";
import { datetimeFormats, numberFormats } from "./formats";
import { DEFAULT_LOCALE } from "./types";
import thMessages from "./locales/th/index";

// Eagerly load English in test mode so Vitest specs work without async load
const messages: Record<string, unknown> = {
  th: thMessages,
};

if (import.meta.env.MODE === "test") {
  // @ts-expect-error: English is lazy-loaded in production but imported synchronously in tests
  const enModule = await import("./locales/en/index");
  messages.en = enModule.default;
}

export const i18n = createI18n({
  legacy: false,
  locale: DEFAULT_LOCALE,
  fallbackLocale: DEFAULT_LOCALE,
  messages,
  datetimeFormats,
  numberFormats,
  missingWarn: !import.meta.env.PROD,
  fallbackWarn: !import.meta.env.PROD,
});

export * from "./types";
export * from "./locale";
export * from "./useLocale";
export * from "./formats";
