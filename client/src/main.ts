import { createApp } from "vue"
import { createPinia } from "pinia"
import { VueQueryPlugin } from "@tanstack/vue-query"

import App from "./App.vue"
import router from "./router"

import { useAuthStore } from "@/stores/auth"

import "./assets/styles/main.css"
import "primeicons/primeicons.css"

async function bootstrap() {
  const app = createApp(App)

  const pinia = createPinia()
  app.use(pinia)

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
