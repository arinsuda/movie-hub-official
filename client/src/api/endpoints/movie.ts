import api from '../index'
import type { PaginatedResult, Movie, MovieDetail, Credits, Video, Genre, TVSeries, TVSeriesDetail } from '@/types'

export const movieApi = {
  // Movies
  getPopular: (page = 1) =>
    api.get<PaginatedResult<Movie>>('/movies/popular', { params: { page } }),

  getNowPlaying: (page = 1) =>
    api.get<PaginatedResult<Movie>>('/movies/now-playing', { params: { page } }),

  getTopRated: (page = 1) =>
    api.get<PaginatedResult<Movie>>('/movies/top-rated', { params: { page } }),

  getUpcoming: (page = 1) =>
    api.get<PaginatedResult<Movie>>('/movies/upcoming', { params: { page } }),

  search: (q: string, page = 1) =>
    api.get<PaginatedResult<Movie>>('/movies/search', { params: { q, page } }),

  getGenres: () =>
    api.get<{ genres: Genre[] }>('/movies/genres'),

  getById: (id: number) =>
    api.get<{ movie: MovieDetail; credits: Credits; videos: Video[] }>(`/movies/${id}`),

  getSimilar: (id: number, page = 1) =>
    api.get<PaginatedResult<Movie>>(`/movies/${id}/similar`, { params: { page } }),

  // TV Series
  getPopularSeries: (page = 1) =>
    api.get<PaginatedResult<TVSeries>>('/tv/popular', { params: { page } }),

  getNowAiringSeries: (page = 1) =>
    api.get<PaginatedResult<TVSeries>>('/tv/now-airing', { params: { page } }),

  getTopRatedSeries: (page = 1) =>
    api.get<PaginatedResult<TVSeries>>('/tv/top-rated', { params: { page } }),

  searchSeries: (q: string, page = 1) =>
    api.get<PaginatedResult<TVSeries>>('/tv/search', { params: { q, page } }),

  getSeriesGenres: () =>
    api.get<{ genres: Genre[] }>('/tv/genres'),

  getSeriesById: (id: number) =>
    api.get<{ series: TVSeriesDetail; credits: Credits; videos: Video[] }>(`/tv/${id}`),

  getSimilarSeries: (id: number, page = 1) =>
    api.get<PaginatedResult<TVSeries>>(`/tv/${id}/similar`, { params: { page } }),
}
