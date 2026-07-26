import { ref } from "vue"
import { toPng } from "html-to-image"
import { useToast } from "@/composables/useToast"
import { useI18n } from "vue-i18n"
import type { MediaType } from "@/types/common"

/** Build API base URL matching the same logic as @/api/index.ts */
function getApiBaseUrl(): string {
  const envUrl = import.meta.env.VITE_API_BASE_URL
  if (!envUrl) return "/api"
  const cleanUrl = (envUrl as string).replace(/\/+$/, "")
  return cleanUrl.endsWith("/api") ? cleanUrl : `${cleanUrl}/api`
}

/** Hosts known to reject browser CORS fetch — skip direct attempt, go straight to proxy */
const CORS_UNFRIENDLY_HOSTS = ["image.tmdb.org"]

function isCorsUnfriendly(url: string): boolean {
  try {
    const host = new URL(url).hostname
    return CORS_UNFRIENDLY_HOSTS.some((h) => host === h || host.endsWith(`.${h}`))
  } catch {
    return false
  }
}

export function useShareImage() {
  let t: (key: string) => string
  try {
    const i18n = useI18n()
    t = i18n.t
  } catch {
    t = (key: string) => key
  }

  const toast = useToast()

  const isExporting = ref(false)
  const blobUrlsToRevoke: string[] = []
  let hasFallback = false

  /** Fetch a remote image as a blob URL to bypass CORS canvas taint */
  async function fetchImageAsBlob(url: string): Promise<string | null> {
    try {
      let response: Response | null = null

      // 1. Skip direct fetch for hosts known to block CORS (e.g. TMDB CDN)
      if (!isCorsUnfriendly(url)) {
        response = await fetch(url, { mode: "cors" }).catch(() => null)
      }

      // 2. If direct fetch failed or was skipped, use backend proxy
      if (!response || !response.ok) {
        const apiBase = getApiBaseUrl()
        const proxyUrl = `${apiBase}/share/proxy-image?url=${encodeURIComponent(url)}`
        response = await fetch(proxyUrl).catch(() => null)
      }

      if (!response || !response.ok) return null

      const blob = await response.blob()
      const blobUrl = URL.createObjectURL(blob)
      blobUrlsToRevoke.push(blobUrl)
      return blobUrl
    } catch {
      return null
    }
  }

  /** Wait for all images inside a node to finish loading */
  async function waitForImages(node: HTMLElement): Promise<void> {
    const images = node.querySelectorAll("img")
    const promises: Promise<void>[] = []
    images.forEach((img) => {
      if (img.complete && img.naturalWidth > 0) return
      promises.push(
        img.decode().catch(() => {
          /* fallback image or broken — ignore */
        }),
      )
    })
    await Promise.all(promises)
  }

  /** Generate a PNG blob from a DOM node (1080×1920) */
  async function generatePng(node: HTMLElement): Promise<Blob | null> {
    if (isExporting.value) return null
    isExporting.value = true
    try {
      await document.fonts.ready
      await waitForImages(node)

      const dataUrl = await toPng(node, {
        pixelRatio: 1,
        width: 1080,
        height: 1920,
        cacheBust: true,
      })

      const response = await fetch(dataUrl)
      return await response.blob()
    } catch {
      toast.error(t("share.failedToGenerate"))
      return null
    } finally {
      isExporting.value = false
    }
  }

  /** Trigger a browser download for a blob */
  function downloadBlob(blob: Blob, filename: string): void {
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = filename
    a.style.display = "none"
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    setTimeout(() => URL.revokeObjectURL(url), 1000)
    toast.success(t("share.downloadSuccess"))
  }

  /** Copy text to clipboard with toast feedback */
  async function copyText(text: string, successKey: string): Promise<void> {
    try {
      await navigator.clipboard.writeText(text)
      toast.success(t(successKey))
    } catch {
      toast.error(t("share.clipboardFailed"))
    }
  }

  /** Create a safe filename from a title */
  function sanitizeFilename(title: string, prefix: string): string {
    const sanitized = title
      .replace(/[^a-zA-Z0-9\u0E00-\u0E7F_-]/g, "_")
      .replace(/_+/g, "_")
      .replace(/^_|_$/g, "")
      .slice(0, 50)
    return `${prefix}_${sanitized || "share"}.png`
  }

  /** Deterministic Unicode-safe review truncation */
  function truncateReview(text: string, maxChars: number = 200): string {
    if (text.length <= maxChars) return text

    // Use Array.from to handle surrogate pairs correctly
    const chars = Array.from(text)
    if (chars.length <= maxChars) return text

    const truncated = chars.slice(0, maxChars).join("")

    // Try to find last space for word boundary (works for English)
    const lastSpace = truncated.lastIndexOf(" ")
    if (lastSpace > maxChars * 0.6) {
      return truncated.slice(0, lastSpace) + "\u2026"
    }

    // For Thai or no-space text, truncate at char limit directly
    return truncated + "\u2026"
  }

  /** Build a caption for poster share */
  function buildPosterCaption(
    title: string,
    year: string | null,
    mediaType: MediaType,
  ): string {
    const typeLabel = t(mediaType === "movie" ? "share.movie" : "share.tvSeries")
    const yearStr = year ? ` (${year})` : ""
    return `\uD83C\uDFAC ${title}${yearStr}\n${typeLabel}\n\n${t("share.sharedFrom")}`
  }

  /** Build a caption for review share */
  function buildReviewCaption(
    title: string,
    authorName: string,
    rating: number,
  ): string {
    return `\u2B50 ${rating.toFixed(1)}/5 \u2014 ${title}\n\uD83D\uDCDD ${t("share.reviewedOn")} \u2022 ${authorName}\n\n${t("share.sharedFrom")}`
  }

  /** Build a share link for the media detail page */
  function buildShareLink(mediaType: MediaType, mediaId: number): string {
    const base = window.location.origin
    const path = mediaType === "movie" ? "movies" : "tv"
    return `${base}/${path}/${mediaId}`
  }

  /** Mark that a fallback was used for an image */
  function markFallbackUsed(): void {
    hasFallback = true
  }

  /** Show fallback warning if any images failed */
  function showFallbackWarningIfNeeded(): void {
    if (hasFallback) {
      toast.warning(t("share.imageFallbackUsed"))
    }
  }

  /** Reset fallback tracking */
  function resetFallbackTracking(): void {
    hasFallback = false
  }

  /** Revoke all blob URLs and reset state */
  function cleanup(): void {
    blobUrlsToRevoke.forEach((url) => URL.revokeObjectURL(url))
    blobUrlsToRevoke.length = 0
    isExporting.value = false
    hasFallback = false
  }

  return {
    isExporting,
    fetchImageAsBlob,
    generatePng,
    downloadBlob,
    copyText,
    sanitizeFilename,
    truncateReview,
    buildPosterCaption,
    buildReviewCaption,
    buildShareLink,
    markFallbackUsed,
    showFallbackWarningIfNeeded,
    resetFallbackTracking,
    cleanup,
  }
}
