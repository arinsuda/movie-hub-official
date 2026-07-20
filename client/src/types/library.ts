export type LibraryTab =
  | "dashboard"
  | "bmol"
  | "watchlist"
  | "watched"
  | "likes"
  | "reviews"

export type LibraryListTab = Exclude<LibraryTab, "dashboard" | "bmol">
