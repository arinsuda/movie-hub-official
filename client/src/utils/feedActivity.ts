import type { Component } from "vue";
import type { ActivityType } from "@/types/feed";
import {
  PenTool,
  MessageSquare,
  Heart,
  Bookmark,
  CheckCircle2,
  Trophy,
  UserPlus,
  HelpCircle
} from "lucide-vue-next";

interface ActivityMeta {
  label: string;
  icon: Component;
}

const ACTIVITY_META: Record<ActivityType, ActivityMeta> = {
  review_created: { label: "รีวิวภาพยนตร์", icon: PenTool },
  review_commented: { label: "แสดงความคิดเห็น", icon: MessageSquare },
  review_liked: { label: "ถูกใจรีวิว", icon: Heart },
  media_liked: { label: "ถูกใจ", icon: Heart },
  watchlist_added: { label: "เพิ่มลงวอทช์ลิสต์", icon: Bookmark },
  watched_added: { label: "ดูแล้ว", icon: CheckCircle2 },
  achievement_unlocked: { label: "ปลดล็อกความสำเร็จ", icon: Trophy },
  user_followed: { label: "ติดตามผู้ใช้", icon: UserPlus },
};

export function getActivityMeta(type: ActivityType): ActivityMeta {
  return ACTIVITY_META[type] ?? { label: type, icon: HelpCircle };
}

// แปลง media_type จาก backend ให้เป็น prefix ของ route ฝั่ง FE
export function mediaRoutePrefix(mediaType: string): "movies" | "tv" {
  return mediaType === "tv" || mediaType === "series" ? "tv" : "movies";
}
