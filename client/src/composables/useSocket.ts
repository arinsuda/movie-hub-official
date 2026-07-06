import { io, type Socket } from "socket.io-client"
import { ref } from "vue"

// ปรับให้ตรงกับ backend ของคุณ — ถ้า BE รันคนละพอร์ต/โดเมนกับ REST API
// ให้ตั้งค่า VITE_SOCKET_URL แยกไว้ใน .env
const SOCKET_URL =
  import.meta.env.VITE_SOCKET_URL ||
  import.meta.env.VITE_API_BASE_URL ||
  window.location.origin

// เก็บ instance ไว้แบบ module-level (singleton) กันสร้างซ้ำเวลา component
// หลายตัวเรียก useSocket() พร้อมกัน
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
  function connect(userId: number): Socket {
    if (socket?.connected || isConnecting.value) return socket as Socket

    isConnecting.value = true

    socket = io(SOCKET_URL, {
      withCredentials: true, // ส่ง cookie (httpOnly session/JWT) ไปด้วย
      transports: ["websocket", "polling"],
      auth: { userId },
      reconnection: true,
      reconnectionAttempts: 10,
      reconnectionDelay: 1500,
      reconnectionDelayMax: 8000,
    })

    socket.on("connect", () => {
      isConnected.value = true
      isConnecting.value = false
      // ให้ server join socket นี้เข้าห้องส่วนตัวของ user
      // (อีกทางคือให้ server อ่าน userId จาก auth handshake แล้ว join ให้เองเลยก็ได้
      // ถ้าทำแบบนั้นตัด emit "join" ทิ้งได้)
      socket?.emit("join", `user:${userId}`)
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
