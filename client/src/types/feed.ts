import type { MediaType } from "@/types";

export type ActivityType =
  | "review_created"
  | "review_commented"
  | "review_liked"
  | "media_liked"
  | "watchlist_added"
  | "watched_added"
  | "achievement_unlocked"
  | "user_followed";

export interface FeedActorSummary {
  id: number;
  username: string;
  display_name: string | null;
  avatar_url: string | null;
}

export interface FeedMediaSummary {
  id: number;
  media_type: MediaType;
  title?: string;
  poster_url?: string | null;
  genres?: string[];
  vote_average?: number;
}

export interface FeedItemResponse {
  id: number;
  type: ActivityType;
  actor: FeedActorSummary;
  media?: FeedMediaSummary;
  review_id?: number;
  comment_id?: number;
  achievement_id?: number;
  library_item_id?: number;
  target_user?: FeedActorSummary;
  message: string;
  visibility: "default" | "public" | "followers" | "private";
  created_at: string;
}

export interface FeedPaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface FeedListResponse {
  items: FeedItemResponse[];
  pagination: FeedPaginationMeta;
}

export interface FeedQueryParams {
  page?: number;
  limit?: number;
}

export interface ActivitySettingsResponse {
  review_created: boolean;
  review_commented: boolean;
  review_liked: boolean;
  media_liked: boolean;
  watchlist_added: boolean;
  watched_added: boolean;
  achievement_unlocked: boolean;
  user_followed: boolean;
}

export type UpdateActivitySettingsRequest = Partial<ActivitySettingsResponse>;
