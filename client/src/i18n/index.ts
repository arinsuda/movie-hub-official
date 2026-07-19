import { createI18n } from 'vue-i18n'
import th from '../locales/th.json'
import en from '../locales/en.json'

export const i18n = createI18n({
  legacy: false,
  locale: 'th',
  fallbackLocale: 'en',
  messages: {
    th,
    en
  }
})
