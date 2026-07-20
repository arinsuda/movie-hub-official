import axios from "axios";
import { useAuthStore } from "@/stores/auth";

const getBaseURL = () => {
  const envUrl = import.meta.env.VITE_API_BASE_URL;
  if (!envUrl) return "/api";
  const cleanUrl = envUrl.replace(/\/+$/, "");
  return cleanUrl.endsWith("/api") ? cleanUrl : `${cleanUrl}/api`;
};

const api = axios.create({
  baseURL: getBaseURL(),
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

let isRefreshing = false;

let failedQueue: {
  resolve: (value?: unknown) => void;
  reject: (error?: any) => void;
}[] = [];

function processQueue(error: any = null) {
  failedQueue.forEach((promise) => {
    if (error) {
      promise.reject(error);
    } else {
      promise.resolve();
    }
  });

  failedQueue = [];
}

api.interceptors.response.use(
  (response) => response,

  async (error) => {
    const originalRequest = error.config;
    const status = error.response?.status;

    // ป้องกัน refresh loop — ไม่ retry ถ้าเป็น request จาก auth routes เหล่านี้
    const isAuthRoute = ["/auth/refresh", "/auth/login", "/auth/logout"].some(
      (path) => originalRequest?.url?.includes(path),
    );
    if (isAuthRoute) {
      return Promise.reject(error);
    }

    // retry เฉพาะ 401 เท่านั้น
    // 404 / 403 / 5xx ไม่ใช่ expired token → reject ออกไปเลย
    if (status !== 401 || originalRequest._retry) {
      return Promise.reject(error);
    }

    originalRequest._retry = true;

    // ถ้ามี refresh request อยู่แล้ว ให้รอใน queue
    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        failedQueue.push({ resolve, reject });
      }).then(() => api(originalRequest));
    }

    isRefreshing = true;

    try {
      await axios.post("/api/auth/refresh", {}, { withCredentials: true });

      processQueue();

      return api(originalRequest);
    } catch (refreshError) {
      processQueue(refreshError);

      const authStore = useAuthStore();
      authStore.clearUser();
      window.location.href = "/login";

      return Promise.reject(refreshError);
    } finally {
      isRefreshing = false;
    }
  },
);

export default api;
