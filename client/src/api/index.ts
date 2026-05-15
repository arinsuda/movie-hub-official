import axios from "axios";
import { useAuthStore } from "@/stores/auth";

const api = axios.create({
  baseURL: "/api",
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

    // กัน refresh loop
    if (originalRequest?.url?.includes("/auth/refresh")) {
      return Promise.reject(error);
    }

    // access token หมด
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      // ถ้ามี refresh request อยู่แล้ว
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({
            resolve,
            reject,
          });
        }).then(() => api(originalRequest));
      }

      isRefreshing = true;

      try {
        // refresh token
        await axios.post(
          "/api/auth/refresh",
          {},
          {
            withCredentials: true,
          },
        );

        processQueue();

        // retry request เดิม
        return api(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError);

        // logout
        const authStore = useAuthStore();

        authStore.clearUser();

        window.location.href = "/login";

        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  },
);

export default api;
