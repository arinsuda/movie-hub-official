import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import { followApi } from "@/api/endpoints/social"
import type { FollowRelationshipStatus } from "@/types/follow"

export function useFollowStats(userId: number) {
  return useQuery({
    queryKey: ["follow-stats", userId],
    queryFn: () => followApi.getFollowStats(userId).then(res => res.data),
    enabled: userId > 0,
  })
}

export function useFollowStatus(userId: number) {
  return useQuery({
    queryKey: ["follow-status", userId],
    queryFn: () => followApi.getFollowStatus(userId).then(res => res.data),
    enabled: userId > 0,
  })
}

export function useFollowActions(userId: number) {
  const queryClient = useQueryClient()

  function invalidate() {
    queryClient.invalidateQueries({ queryKey: ["follow-status", userId] })
    queryClient.invalidateQueries({ queryKey: ["follow-stats", userId] })
  }

  const follow = useMutation({
    mutationFn: () => followApi.follow(userId),
    onSuccess: res => {
      queryClient.setQueryData<FollowRelationshipStatus>(
        ["follow-status", userId],
        prev => ({
          is_following: res.data.status === "accepted",
          follow_status: res.data.status,
          is_followed_by: prev?.is_followed_by ?? false,
        }),
      )
      invalidate()
    },
  })

  const unfollow = useMutation({
    mutationFn: () => followApi.unfollow(userId),
    onSuccess: () => {
      queryClient.setQueryData<FollowRelationshipStatus>(
        ["follow-status", userId],
        prev => ({
          is_following: false,
          follow_status: undefined,
          is_followed_by: prev?.is_followed_by ?? false,
        }),
      )
      invalidate()
    },
  })

  return { follow, unfollow }
}
