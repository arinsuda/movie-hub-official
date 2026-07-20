export const datetimeFormats = {
  th: {
    short: { year: "numeric", month: "short", day: "numeric" },
    medium: { year: "numeric", month: "long", day: "numeric" },
    full: { year: "numeric", month: "long", day: "numeric", weekday: "long" },
    year: { year: "numeric" },
  },
  en: {
    short: { year: "numeric", month: "short", day: "numeric" },
    medium: { year: "numeric", month: "long", day: "numeric" },
    full: { year: "numeric", month: "long", day: "numeric", weekday: "long" },
    year: { year: "numeric" },
  },
} as const;

export const numberFormats = {
  th: {
    compact: { notation: "compact", compactDisplay: "short" },
    decimal: { style: "decimal", minimumFractionDigits: 1, maximumFractionDigits: 1 },
    percent: { style: "percent", minimumFractionDigits: 0, maximumFractionDigits: 1 },
  },
  en: {
    compact: { notation: "compact", compactDisplay: "short" },
    decimal: { style: "decimal", minimumFractionDigits: 1, maximumFractionDigits: 1 },
    percent: { style: "percent", minimumFractionDigits: 0, maximumFractionDigits: 1 },
  },
} as const;

export function formatRelativeTime(dateInput: string | Date | number, locale: string): string {
  const date = typeof dateInput === "string" || typeof dateInput === "number" ? new Date(dateInput) : dateInput;
  if (isNaN(date.getTime())) return "";

  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSecs = Math.floor(diffMs / 1000);
  const diffMins = Math.floor(diffSecs / 60);
  const diffHours = Math.floor(diffMins / 60);
  const diffDays = Math.floor(diffHours / 24);

  const isTh = locale === "th";

  if (diffSecs < 10) {
    return isTh ? "เมื่อครู่" : "just now";
  }
  if (diffSecs < 60) {
    return isTh ? `${diffSecs} วินาทีที่แล้ว` : `${diffSecs} seconds ago`;
  }
  if (diffMins < 60) {
    return isTh ? `${diffMins} นาทีที่แล้ว` : `${diffMins} ${diffMins === 1 ? "minute" : "minutes"} ago`;
  }
  if (diffHours < 24) {
    return isTh ? `${diffHours} ชั่วโมงที่แล้ว` : `${diffHours} ${diffHours === 1 ? "hour" : "hours"} ago`;
  }
  if (diffDays < 7) {
    return isTh ? `${diffDays} วันที่แล้ว` : `${diffDays} ${diffDays === 1 ? "day" : "days"} ago`;
  }

  const options: Intl.DateTimeFormatOptions = { year: "numeric", month: "short", day: "numeric" };
  return new Intl.DateTimeFormat(locale, options).format(date);
}
