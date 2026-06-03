package movie_module

type MediaType string
type ListType string

const (
	MediaMovie  MediaType = "movie"
	MediaSeries MediaType = "tv"
)

const (
	ListWatchlist ListType = "watchlist"
	ListFavorite  ListType = "favorite"
	ListWatched   ListType = "watched"
)
