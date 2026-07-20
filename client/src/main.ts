import { createApp } from "vue"
import { createPinia } from "pinia"
import { VueQueryPlugin } from "@tanstack/vue-query"

import App from "./App.vue"
import router from "./router"

import { useAuthStore } from "@/stores/auth"
import { VueDatePicker } from "@vuepic/vue-datepicker"
import "@vuepic/vue-datepicker/dist/main.css"
import "./assets/styles/main.css"
import "primeicons/primeicons.css"
import PrimeVue from "primevue/config"
import { i18n, getSavedLocale, setLocale } from "./i18n"

async function bootstrap() {
  const app = createApp(App)

  // set saved or default locale
  const savedLocale = getSavedLocale()
  await setLocale(savedLocale)

  const pinia = createPinia()
  app.use(pinia)
  app.use(PrimeVue, { ripple: true })
  app.use(i18n)
  app.component("VueDatePicker", VueDatePicker)

  // restore session
  const authStore = useAuthStore()
  await authStore.fetchMe()

  app.use(router)

  app.use(VueQueryPlugin, {
    queryClientConfig: {
      defaultOptions: {
        queries: {
          staleTime: 1000 * 60 * 5,
          retry: 1,
        },
      },
    },
  })

  app.mount("#app")
}

bootstrap()
