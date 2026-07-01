
import api from "../index"

import type {
  AchievementQueryParams,
  ListAchievementsResponse,
  ListUserAchievementsResponse,
  UserAchievementQueryParams,
} from "@/types/achievement"

export const achievementApi = {
  listAll(params: AchievementQueryParams = {}) {
    return api
      .get<ListAchievementsResponse>("/achievements", { params })
      .then(res => res.data)
  },

  listByUser(userId: number, params: UserAchievementQueryParams = {}) {
    return api
      .get<ListUserAchievementsResponse>(`/users/${userId}/achievements`, {
        params,
      })
      .then(res => res.data)
  },
}
