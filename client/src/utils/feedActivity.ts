import type { ActivityType } from "@/types/feed";

interface ActivityMeta {
  label: string;
  icon: string;
}

const ACTIVITY_META: Record<ActivityType, ActivityMeta> = {
  review_created: { label: "รีวิว", icon: "📝" },
  review_commented: { label: "แสดงความคิดเห็น", icon: "💬" },
  review_liked: { label: "ถูกใจรีวิว", icon: "❤️" },
  media_liked: { label: "ถูกใจ", icon: "❤️" },
  watchlist_added: { label: "เพิ่มลงวอทช์ลิสต์", icon: "🔖" },
  watched_added: { label: "ดูแล้ว", icon: "✅" },
  achievement_unlocked: { label: "ปลดล็อกความสำเร็จ", icon: "🏆" },
};

export function getActivityMeta(type: ActivityType): ActivityMeta {
  return ACTIVITY_META[type] ?? { label: type, icon: "•" };
}

// แปลง media_type จาก backend ให้เป็น prefix ของ route ฝั่ง FE
// ปรับ mapping ตรงนี้ให้ตรงกับ router จริง (เช่นถ้าใช้ /movie /series แทน /movies /tv)
export function mediaRoutePrefix(mediaType: string): "movies" | "tv" {
  return mediaType === "tv" || mediaType === "series" ? "tv" : "movies";
}
