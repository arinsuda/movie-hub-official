import { useRouter } from "vue-router"

export function useProfileNav() {
  const router = useRouter()

  function goToProfile(userId: number) {
    router.push({ name: "user-profile", params: { userId: userId } })
  }

  return { goToProfile }
}
