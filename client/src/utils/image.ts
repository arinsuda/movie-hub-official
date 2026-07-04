// client/src/utils/image.ts

const TMDB_IMAGE_BASE = "https://image.tmdb.org/t/p";

export type TmdbImageSize =
  | "w92"
  | "w154"
  | "w185"
  | "w342"
  | "w500"
  | "w780"
  | "original";

/**
 * สร้าง URL เต็มสำหรับรูปจาก TMDB
 * @param path poster_path หรือ backdrop_path ที่ได้จาก API (เช่น "/abc123.jpg")
 * @param size ขนาดรูปที่ต้องการ (default: w342)
 * @returns URL เต็ม หรือ null ถ้าไม่มี path
 */
export function getTmdbImageUrl(
  path: string | null | undefined,
  size: TmdbImageSize = "w342",
): string | null {
  if (!path) return null;
  return `${TMDB_IMAGE_BASE}/${size}${path}`;
}

/**
 * เหมือน getTmdbImageUrl แต่ return placeholder แทน null
 * ใช้ตรงๆ กับ :src ได้เลยโดยไม่ต้องเช็ค null เอง
 */
export function getTmdbImageUrlOrPlaceholder(
  path: string | null | undefined,
  size: TmdbImageSize = "w342",
): string {
  return getTmdbImageUrl(path, size) ?? "/placeholder.jpg";
}
