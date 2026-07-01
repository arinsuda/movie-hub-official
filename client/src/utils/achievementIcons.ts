import {
  UserCheck,
  PenLine,
  Eye,
  Heart,
  Flame,
  Star,
  Trophy,
  type LucideIcon,
} from "lucide-vue-next"

const ICON_MAP: Record<string, LucideIcon> = {
  profile_complete: UserCheck,
  review_count: PenLine,
  watchlist_count: Eye,
  like_count: Heart,
  login_streak: Flame,
  rating_count: Star,
}

export function getAchievementIcon(actionType: string): LucideIcon {
  return ICON_MAP[actionType] ?? Trophy
}
