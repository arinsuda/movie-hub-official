export interface ExtractedColor {
  r: number
  g: number
  b: number
  hex: string
  rgbStr: string
}

const DEFAULT_COLOR: ExtractedColor = {
  r: 229,
  g: 9,
  b: 20,
  hex: "#e50914",
  rgbStr: "rgb(229, 9, 20)",
}

/**
 * Extract dominant vibrant color from an image URL (preferably a blob URL to avoid CORS).
 */
export async function extractDominantColor(
  imageSrc: string | null | undefined
): Promise<ExtractedColor> {
  if (!imageSrc) return DEFAULT_COLOR

  return new Promise((resolve) => {
    const img = new Image()
    img.crossOrigin = "Anonymous"

    img.onload = () => {
      try {
        const canvas = document.createElement("canvas")
        const ctx = canvas.getContext("2d")
        if (!ctx) return resolve(DEFAULT_COLOR)

        const width = 60
        const height = 90
        canvas.width = width
        canvas.height = height

        ctx.drawImage(img, 0, 0, width, height)

        const imageData = ctx.getImageData(0, 0, width, height)
        const data = imageData.data

        let totalR = 0
        let totalG = 0
        let totalB = 0
        let count = 0

        // Find vibrant pixels (not too dark, not too bright, reasonable saturation)
        let vibrantR = 0
        let vibrantG = 0
        let vibrantB = 0
        let maxScore = -1

        for (let i = 0; i < data.length; i += 4) {
          const r = data[i] ?? 0
          const g = data[i + 1] ?? 0
          const b = data[i + 2] ?? 0
          const a = data[i + 3] ?? 255

          if (a < 128) continue // ignore transparent pixels

          const max = Math.max(r, g, b)
          const min = Math.min(r, g, b)
          const brightness = (r + g + b) / 3
          const saturation = max === 0 ? 0 : (max - min) / max

          // Skip extreme black or white
          if (brightness < 20 || brightness > 240) continue

          totalR += r
          totalG += g
          totalB += b
          count++

          // Score pixel based on saturation and moderate brightness
          const score = saturation * 2 + (1 - Math.abs(brightness - 128) / 128)
          if (score > maxScore) {
            maxScore = score
            vibrantR = r
            vibrantG = g
            vibrantB = b
          }
        }

        if (count === 0 && maxScore < 0) {
          return resolve(DEFAULT_COLOR)
        }

        // Use vibrant color if good score, else average
        const r = maxScore > 0.5 ? vibrantR : Math.round(totalR / (count || 1))
        const g = maxScore > 0.5 ? vibrantG : Math.round(totalG / (count || 1))
        const b = maxScore > 0.5 ? vibrantB : Math.round(totalB / (count || 1))

        const hex = `#${((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)}`
        resolve({
          r,
          g,
          b,
          hex,
          rgbStr: `rgb(${r}, ${g}, ${b})`,
        })
      } catch {
        resolve(DEFAULT_COLOR)
      }
    }

    img.onerror = () => {
      resolve(DEFAULT_COLOR)
    }

    img.src = imageSrc
  })
}
