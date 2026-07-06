import { io, type Socket } from "socket.io-client"
import { ref } from "vue"

const SOCKET_URL =
  import.meta.env.VITE_SOCKET_URL ||
  import.meta.env.VITE_API_BASE_URL ||
  window.location.origin

let socket: Socket | null = null
const isConnected = ref(false)
const isConnecting = ref(false)

export function useSocket() {
  /**
   * เปิดการเชื่อมต่อ (no-op ถ้าเชื่อมอยู่แล้ว)
   * userId ใช้ยืนยันตัวตนฝั่ง server เพื่อ join ห้องเฉพาะ user นั้น
   * ถ้า backend ใช้ cookie/JWT อยู่แล้วผ่าน withCredentials ก็ยังส่ง userId
   * เป็น fallback ไว้เผื่อ server อยากใช้ยืนยันห้องอีกชั้น
   */
  function connect(accessToken: string): Socket {
    if (socket?.connected || isConnecting.value) return socket as Socket

    isConnecting.value = true

    socket = io(SOCKET_URL, {
      withCredentials: true,
      transports: ["websocket", "polling"],
      auth: { token: accessToken },
      reconnection: true,
      reconnectionAttempts: 10,
      reconnectionDelay: 1500,
      reconnectionDelayMax: 8000,
    })

    socket.on("connect", () => {
      isConnected.value = true
      isConnecting.value = false

      socket?.emit("join", `user:${accessToken}`)
    })

    socket.on("disconnect", reason => {
      isConnected.value = false
      if (import.meta.env.DEV) {
        console.warn("[socket] disconnected:", reason)
      }
    })

    socket.on("connect_error", err => {
      isConnecting.value = false
      if (import.meta.env.DEV) {
        console.error("[socket] connect_error:", err.message)
      }
    })

    return socket
  }

  function disconnect() {
    socket?.removeAllListeners()
    socket?.disconnect()
    socket = null
    isConnected.value = false
    isConnecting.value = false
  }

  function on<T = unknown>(event: string, handler: (payload: T) => void) {
    socket?.on(event, handler as (...args: unknown[]) => void)
  }

  function off(event: string, handler?: (...args: unknown[]) => void) {
    socket?.off(event, handler)
  }

  function emit(event: string, payload?: unknown) {
    if (!socket?.connected) {
      if (import.meta.env.DEV) {
        console.warn(`[socket] tried to emit "${event}" while disconnected`)
      }
      return
    }
    socket.emit(event, payload)
  }

  return { connect, disconnect, on, off, emit, isConnected, isConnecting }
}
