import api from "../index";
import type {
  ReviewResponse,
  CreateReviewRequest,
  UpdateReviewRequest,
  CommentResponse,
  CreateCommentRequest,
  LibraryItemResponse,
  AddItemRequest,
  MediaStatusResponse,
  PaginatedResponse,
  FollowStatsResponse,
  FollowUserResponse,
  FeedItemResponse,
  MediaType,
  ListType,
} from "@/types";

// ── Review ─────────────────────────────────────────────────────────
export const reviewApi = {
  createReview: (userId: number, data: CreateReviewRequest) =>
    api.post<{ review: ReviewResponse }>(`/users/${userId}/reviews`, data),

  getUserReviews: (userId: number, page = 1, limit = 20) =>
    api.get<PaginatedResponse<ReviewResponse>>(`/users/${userId}/reviews`, {
      params: { page, limit },
    }),

  getMediaReviews: (
    mediaType: "movies" | "series",
    mediaId: number,
    params?: { page?: number; limit?: number; sort?: string },
  ) =>
    api.get<PaginatedResponse<ReviewResponse>>(
      `/${mediaType}/${mediaId}/reviews`,
      {
        params: { page: 1, limit: 20, ...params },
      },
    ),

  updateReview: (userId: number, reviewId: number, data: UpdateReviewRequest) =>
    api.patch<{ review: ReviewResponse }>(
      `/users/${userId}/reviews/${reviewId}`,
      data,
    ),

  deleteReview: (userId: number, reviewId: number) =>
    api.delete(`/users/${userId}/reviews/${reviewId}`),

  likeReview: (reviewId: number) => api.post(`/reviews/${reviewId}/likes`),

  unlikeReview: (reviewId: number) => api.delete(`/reviews/${reviewId}/likes`),

  getComments: (reviewId: number, page = 1, limit = 20) =>
    api.get<PaginatedResponse<CommentResponse>>(
      `/reviews/${reviewId}/comments`,
      {
        params: { page, limit },
      },
    ),

  createComment: (reviewId: number, data: CreateCommentRequest) =>
    api.post<{ comment: CommentResponse }>(
      `/reviews/${reviewId}/comments`,
      data,
    ),

  updateComment: (reviewId: number, commentId: number, body: string) =>
    api.patch<{ comment: CommentResponse }>(
      `/reviews/${reviewId}/comments/${commentId}`,
      { body },
    ),

  deleteComment: (reviewId: number, commentId: number) =>
    api.delete(`/reviews/${reviewId}/comments/${commentId}`),
};

// ── Library ────────────────────────────────────────────────────────
export const libraryApi = {
  addItem: (userId: number, data: AddItemRequest) =>
    api.post<{ item: LibraryItemResponse }>(`/users/${userId}/library`, data),

  getLibrary: (
    userId: number,
    params?: { list_type?: ListType; media_type?: MediaType },
  ) =>
    api.get<{ items: LibraryItemResponse[] }>(`/users/${userId}/library`, {
      params,
    }),

  getMediaStatus: (userId: number, mediaId: number, mediaType: MediaType) =>
    api.get<MediaStatusResponse>(`/users/${userId}/library/status`, {
      params: { media_id: mediaId, media_type: mediaType },
    }),

  updateItem: (
    userId: number,
    itemId: number,
    data: { watched_at?: string; tags?: string[]; note?: string },
  ) =>
    api.patch<{ item: LibraryItemResponse }>(
      `/users/${userId}/library/${itemId}`,
      data,
    ),

  removeItem: (userId: number, itemId: number) =>
    api.delete(`/users/${userId}/library/${itemId}`),
};

// ── Follow / Feed ──────────────────────────────────────────────────
export const followApi = {
  follow: (userId: number) => api.post(`/users/${userId}/follow`),

  unfollow: (userId: number) => api.delete(`/users/${userId}/follow`),

  getFollowStats: (userId: number) =>
    api.get<FollowStatsResponse>(`/users/${userId}/follow-stats`),

  getFollowers: (userId: number, page = 1, limit = 20) =>
    api.get<PaginatedResponse<FollowUserResponse>>(
      `/users/${userId}/followers`,
      {
        params: { page, limit },
      },
    ),

  getFollowing: (userId: number, page = 1, limit = 20) =>
    api.get<PaginatedResponse<FollowUserResponse>>(
      `/users/${userId}/following`,
      {
        params: { page, limit },
      },
    ),

  getFeed: (page = 1, limit = 20) =>
    api.get<PaginatedResponse<FeedItemResponse>>("/feed", {
      params: { page, limit },
    }),
};

// ── User ───────────────────────────────────────────────────────────
export const userApi = {
  getProfile: (userId: number) => api.get(`/users/${userId}`),

  updateProfile: (userId: number, data: object) =>
    api.patch(`/users/${userId}`, data),

  deleteUser: (userId: number) => api.delete(`/users/${userId}`),

  updateFavoriteGenres: (userId: number, genres: number[]) =>
    api.patch(`/users/${userId}/genres`, { favorite_genres: genres }),
};
