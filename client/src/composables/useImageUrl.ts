const TMDB_IMAGE_BASE = 'https://image.tmdb.org/t/p'

type ImageSize = 'w92' | 'w154' | 'w185' | 'w342' | 'w500' | 'w780' | 'original'

export function useImageUrl() {
  function tmdbImage(path: string | null | undefined, size: ImageSize = 'w500'): string {
    if (!path) return '/placeholder-poster.png'
    return `${TMDB_IMAGE_BASE}/${size}${path}`
  }

  function backdropImage(path: string | null | undefined): string {
    return tmdbImage(path, 'w780')
  }

  function posterImage(path: string | null | undefined): string {
    return tmdbImage(path, 'w342')
  }

  return { tmdbImage, backdropImage, posterImage }
}
