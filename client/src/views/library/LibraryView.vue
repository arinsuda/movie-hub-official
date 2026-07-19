<template>
  <div class="library-root">
    <div class="grain" aria-hidden="true" />

    <!-- Loading State -->
    <div v-if="loading" class="library-loading">
      <div class="loading-ring" />
      <span class="loading-text">{{ $t("common.loading") }}</span>
    </div>

    <!-- Main Layout -->
    <div v-else class="library-layout">
      <!-- Back to Profile Link -->
      <router-link :to="`/users/${userId}`" class="back-link">
        <ChevronLeft :size="16" />
        <span>{{ $t("navigation.profile") }}</span>
      </router-link>

      <!-- Profile Header Summary Card -->
      <header class="profile-summary-card">
        <div class="user-info-section">
          <div class="avatar-wrapper">
            <div class="avatar-ring">
              <img
                v-if="profileUser?.avatar_url"
                :src="profileUser.avatar_url"
                :alt="profileUser.username"
                class="avatar-img"
              />
              <div v-else class="avatar-fallback">
                <UserIcon :size="28" />
              </div>
            </div>
            <span class="level-chip">Lv.{{ profileUser?.level ?? 1 }}</span>
          </div>
          <div class="user-meta">
            <h1 class="user-name">{{ profileUser?.display_name || profileUser?.username }}</h1>
            <p class="user-tag">@{{ profileUser?.username }}</p>
            <p class="user-bio">{{ profileUser?.bio || $t("profile.noBio") }}</p>
          </div>
        </div>
        <div class="user-counts-section">
          <div class="count-box">
            <span class="count-val">{{ profileUser?.review_count ?? 0 }}</span>
            <span class="count-lbl">{{ $t("library.tabs.reviews") }}</span>
          </div>
          <div class="count-box">
            <span class="count-val">{{ profileUser?.follower_count ?? 0 }}</span>
            <span class="count-lbl">{{ $t("profile.stats.followers") }}</span>
          </div>
          <div class="count-box">
            <span class="count-val">{{ profileUser?.following_count ?? 0 }}</span>
            <span class="count-lbl">{{ $t("profile.stats.following") }}</span>
          </div>
        </div>
      </header>

      <!-- Navigation Tabs -->
      <nav class="library-tabs">
        <button
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === 'dashboard' }"
          @click="activeTab = 'dashboard'"
        >
          <Flame :size="16" />
          <span>{{ $t("library.dashboard.title") }}</span>
        </button>
        <button
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === 'bmol' }"
          @click="activeTab = 'bmol'"
        >
          <Trophy :size="16" />
          <span>{{ $t("library.bmol.title") }}</span>
          <span class="tab-count">{{ bmolItems.length }}</span>
        </button>
        <button
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === 'watchlist' }"
          @click="activeTab = 'watchlist'"
        >
          <Bookmark :size="16" />
          <span>{{ $t("library.tabs.watchlist") }}</span>
          <span class="tab-count">{{ watchlistItems.length }}</span>
        </button>
        <button
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === 'watched' }"
          @click="activeTab = 'watched'"
        >
          <Eye :size="16" />
          <span>{{ $t("library.tabs.watched") }}</span>
          <span class="tab-count">{{ watchedItems.length }}</span>
        </button>
        <button
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === 'likes' }"
          @click="activeTab = 'likes'"
        >
          <Heart :size="16" />
          <span>{{ $t("library.tabs.likes") }}</span>
          <span class="tab-count">{{ likedItems.length }}</span>
        </button>
        <button
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === 'reviews' }"
          @click="activeTab = 'reviews'"
        >
          <Star :size="16" />
          <span>{{ $t("library.tabs.reviews") }}</span>
          <span class="tab-count">{{ userReviews.length }}</span>
        </button>
      </nav>

      <!-- TAB CONTENT: Dashboard -->
      <section v-if="activeTab === 'dashboard'" class="dashboard-section fade-in-up">
        <!-- Highlight stats grid -->
        <div class="stats-grid">
          <!-- Total Watched -->
          <div class="stat-card">
            <div class="stat-icon-box watched-theme">
              <Eye :size="22" />
            </div>
            <div class="stat-data">
              <h3>{{ $t("library.dashboard.totalWatched") }}</h3>
              <p class="stat-value">{{ watchedItems.length }}</p>
            </div>
          </div>

          <!-- Total Reviews & Trend -->
          <div class="stat-card">
            <div class="stat-icon-box review-theme">
              <Star :size="22" />
            </div>
            <div class="stat-data">
              <h3>{{ $t("library.dashboard.totalReviews") }}</h3>
              <p class="stat-value">{{ userReviews.length }}</p>
              <div v-if="userReviews.length > 0" class="trend-badge" :class="trendClass">
                <component :is="trendIcon" :size="12" />
                <span class="trend-text">
                  {{
                    reviewsStats.diff > 0
                      ? $t("library.dashboard.moreThanLastMonth", { count: reviewsStats.diff })
                      : reviewsStats.diff < 0
                      ? $t("library.dashboard.lessThanLastMonth", { count: Math.abs(reviewsStats.diff) })
                      : $t("library.dashboard.equalLastMonth")
                  }}
                </span>
              </div>
            </div>
          </div>

          <!-- Average Rating -->
          <div class="stat-card">
            <div class="stat-icon-box rating-theme">
              <Star :size="22" />
            </div>
            <div class="stat-data">
              <h3>{{ $t("library.dashboard.averageRating") }}</h3>
              <div class="rating-display">
                <span class="stat-value">{{ averageRating > 0 ? averageRating : "—" }}</span>
                <span class="max-val">/ 5.0</span>
              </div>
              <div class="stars-row">
                <span
                  v-for="s in 5"
                  :key="s"
                  class="star-dot"
                  :class="{ 'star-dot--active': s <= Math.round(averageRating) }"
                >★</span>
              </div>
            </div>
          </div>

          <!-- Dynamic Top Genre -->
          <div class="stat-card">
            <div class="stat-icon-box genre-theme">
              <Trophy :size="22" />
            </div>
            <div class="stat-data">
              <h3>{{ $t("library.dashboard.topGenre") }}</h3>
              <p class="stat-value genre-text">{{ topGenre || $t("library.dashboard.noWatchedData") }}</p>
            </div>
          </div>
        </div>

        <!-- Charts Section (Rating Distribution + Genres Breakdown) -->
        <div class="charts-row">
          <!-- Rating Distribution CSS Bar Chart -->
          <div class="chart-card">
            <h4 class="chart-title">{{ $t("library.dashboard.averageRating") }}</h4>
            <div class="distribution-list">
              <div
                v-for="dist in ratingDistribution"
                :key="dist.rating"
                class="dist-row"
              >
                <span class="dist-label">{{ dist.rating.toFixed(1) }} ★</span>
                <div class="dist-bar-wrapper">
                  <div
                    class="dist-bar-fill"
                    :style="{ width: `${dist.percentage}%` }"
                    aria-hidden="true"
                  />
                </div>
                <span class="dist-count">{{ dist.count }}</span>
              </div>
            </div>
          </div>

          <!-- Top Genres Frequency List -->
          <div class="chart-card">
            <h4 class="chart-title">{{ $t("library.dashboard.topGenre") }}</h4>
            <div class="genres-ranking">
              <div
                v-for="(genre, index) in topGenresList"
                :key="genre.name"
                class="genre-rank-row"
              >
                <div class="rank-badge" :class="`rank-${index + 1}`">{{ index + 1 }}</div>
                <span class="genre-rank-name">{{ genre.name }}</span>
                <span class="genre-rank-count">{{ genre.count }} {{ $t("library.tabs.watched") }}</span>
              </div>
              <div v-if="topGenresList.length === 0" class="genres-empty">
                {{ $t("library.dashboard.noWatchedData") }}
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- TAB CONTENT: BMOL (Best Movie/TV of Life) -->
      <section v-else-if="activeTab === 'bmol'" class="bmol-section fade-in-up">
        <!-- Media Type Navigation Subtabs (For View Mode) -->
        <div v-if="selectedRankDetail === null" class="bmol-view-header">
          <h2 class="bmol-section-title">
            {{ bmolSubTab === 'movie' ? $t("library.bmol.movieTitle") : $t("library.bmol.tvTitle") }}
          </h2>
          <div class="bmol-view-selector">
            <button
              class="view-selector-btn"
              :class="{ 'view-selector-btn--active': bmolSubTab === 'movie' }"
              @click="bmolSubTab = 'movie'"
            >
              {{ $t("library.filters.movies") }}
            </button>
            <button
              class="view-selector-btn"
              :class="{ 'view-selector-btn--active': bmolSubTab === 'tv' }"
              @click="bmolSubTab = 'tv'"
            >
              {{ $t("library.filters.tvSeries") }}
            </button>
          </div>
        </div>

        <!-- BREADCRUMB HEADER (When viewing a specific Rank's detail list) -->
        <div v-if="selectedRankDetail !== null" class="bmol-breadcrumb-card">
          <div class="breadcrumb-left">
            <button class="btn-back-circle" :title="$t('library.bmol.backToList')" @click="goBackToList">
              <ChevronLeft :size="16" />
            </button>
            <div class="breadcrumb-path">
              <span class="path-link" @click="goBackToList">{{ $t("library.bmol.title") }}</span>
              <span class="path-separator">/</span>
              <span class="path-current">
                <span class="rank-badge-glow">
                  {{ $t("library.bmol.rank", { rank: selectedRankDetail }) }}
                </span>
              </span>
            </div>
          </div>
          <div class="breadcrumb-right">
            <span class="detail-count-chip">
              {{ bmolFilteredItems.length }} {{ bmolFilteredItems.length === 1 ? 'item' : 'items' }}
            </span>
          </div>
        </div>

        <!-- DETAIL VIEW: When a specific Rank is selected -->
        <div v-if="selectedRankDetail !== null" class="bmol-detail-view fade-in-up">
          <!-- Detail Toolbar: Search + Actions -->
          <div class="detail-toolbar">
            <div class="detail-search-box">
              <Search :size="15" class="detail-search-icon" />
              <input
                v-model="rankFilterQuery"
                type="text"
                :placeholder="$t('library.bmol.filterPlaceholder')"
                class="detail-search-input"
              />
              <button
                v-if="rankFilterQuery.length > 0"
                class="detail-search-clear"
                @click="rankFilterQuery = ''"
              >
                ✕
              </button>
            </div>
            <div class="detail-toolbar-actions">
              <span class="detail-result-count">
                {{ bmolFilteredItems.length }}
                {{ bmolFilteredItems.length === 1 ? 'item' : 'items' }}
              </span>
              <button
                v-if="isOwner"
                class="btn-quick-add"
                @click="openSpotlight(selectedRankDetail!)"
              >
                + {{ $t("library.bmol.addMedia") }}
              </button>
            </div>
          </div>

          <!-- Ranks List inside Rank Detail -->
          <div v-if="bmolFilteredItems.length > 0" class="bmol-detail-grid">
            <div
              v-for="item in bmolFilteredItems"
              :key="item.id"
              class="bmol-item-card"
            >
              <router-link
                :to="item.media_type === 'tv' ? `/tv/${item.media.id}` : `/movies/${item.media.id}`"
                class="bmol-poster-frame"
              >
                <img
                  v-if="item.media.poster_url"
                  :src="`${TMDB_IMG}${item.media.poster_url}`"
                  :alt="item.media.title"
                  loading="lazy"
                />
                <div class="poster-fallback" v-else>
                  <Film :size="18" />
                </div>
              </router-link>
              <div class="bmol-meta">
                <router-link
                  :to="item.media_type === 'tv' ? `/tv/${item.media.id}` : `/movies/${item.media.id}`"
                  class="bmol-title"
                >
                  {{ item.media.title }}
                </router-link>
                <span class="bmol-rating">★ {{ item.media.vote_average.toFixed(1) }}</span>

                <!-- Actions (Only for Owner) -->
                <div v-if="isOwner" class="bmol-actions">
                  <button
                    class="bmol-action-btn"
                    :title="$t('library.bmol.increaseRank')"
                    :disabled="item.rank <= 1"
                    @click="increaseRank(item)"
                  >
                    ▲
                  </button>
                  <button
                    class="bmol-action-btn"
                    :title="$t('library.bmol.decreaseRank')"
                    @click="decreaseRank(item)"
                  >
                    ▼
                  </button>
                  <button
                    class="bmol-action-btn bmol-action-btn--danger"
                    @click="triggerRemoveBmolItem(item)"
                  >
                    ✕
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="empty-state-card">
            <p>{{ $t("library.bmol.errorEmptySearch") }}</p>
          </div>
        </div>

        <!-- DEFAULT VIEW: Grouped Ranks Grid List -->
        <div v-else>
          <!-- ADD FORM: Search bar to add media to any rank (Owner only) -->
          <div v-if="isOwner" class="bmol-add-card">
            <h3 class="bmol-add-title">
              {{ bmolSubTab === 'movie' ? $t("library.bmol.addMovie") : $t("library.bmol.addTV") }}
            </h3>
            <div class="bmol-add-form">
              <!-- Search Input -->
              <div v-if="!bmolAddSelected" class="bmol-add-search-wrapper">
                <Search :size="15" class="bmol-add-search-icon" />
                <input
                  v-model="bmolAddQuery"
                  type="text"
                  :placeholder="$t('library.bmol.searchPlaceholder')"
                  class="bmol-add-search-input"
                />
                <div v-if="bmolAddLoading" class="bmol-add-spinner">
                  <div class="loading-ring small-ring" />
                </div>

                <!-- Search Results Dropdown -->
                <div v-if="bmolAddResults.length > 0" class="bmol-add-dropdown">
                  <div
                    v-for="res in bmolAddResults"
                    :key="res.id"
                    class="bmol-add-dropdown-item"
                    @click="selectBmolAddResult(res)"
                  >
                    <img
                      v-if="res.poster_path"
                      :src="`${TMDB_IMG}${res.poster_path}`"
                      :alt="getMediaTitle(res)"
                      class="bmol-add-dropdown-poster"
                    />
                    <div v-else class="bmol-add-dropdown-poster-fallback">
                      <Film :size="12" />
                    </div>
                    <div class="bmol-add-dropdown-info">
                      <p class="bmol-add-dropdown-title">{{ getMediaTitle(res) }}</p>
                      <p class="bmol-add-dropdown-date">{{ getMediaDate(res) || '—' }}</p>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Selected Preview -->
              <div v-if="bmolAddSelected" class="bmol-add-selected">
                <div class="bmol-add-selected-info">
                  <img
                    v-if="bmolAddSelected.poster_path"
                    :src="`${TMDB_IMG}${bmolAddSelected.poster_path}`"
                    :alt="getMediaTitle(bmolAddSelected)"
                    class="bmol-add-selected-poster"
                  />
                  <div class="bmol-add-selected-detail">
                    <p class="bmol-add-selected-title">{{ getMediaTitle(bmolAddSelected) }}</p>
                    <p class="bmol-add-selected-date">{{ getMediaDate(bmolAddSelected) || '—' }}</p>
                  </div>
                  <button class="bmol-add-clear-btn" @click="clearBmolAddSelection">✕</button>
                </div>
                <div class="bmol-add-rank-row">
                  <label class="bmol-add-rank-label">{{ $t("library.bmol.rank", { rank: '' }).replace('#', '') }}</label>
                  <input
                    v-model.number="bmolAddRank"
                    type="number"
                    min="1"
                    class="bmol-add-rank-input"
                  />
                  <button class="bmol-add-submit-btn" @click="submitBmolAdd">
                    {{ $t("library.bmol.addMedia") }}
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Paginated Ranks -->
          <div v-if="groupedBmolItems.length > 0" class="bmol-ranks-list">
            <div
              v-for="group in paginatedBmolGroups"
              :key="group.rank"
              class="bmol-rank-group"
            >
              <div class="bmol-rank-header">
                <span class="rank-number-tag">{{ $t("library.bmol.rank", { rank: group.rank }) }}</span>
                <div class="rank-line" />
                <!-- Quick Add Button (Only for Owner) -->
                <button
                  v-if="isOwner"
                  class="btn-quick-add"
                  :title="$t('library.bmol.addMedia')"
                  @click="openSpotlight(group.rank)"
                >
                  + {{ $t("library.bmol.addMedia") }}
                </button>
              </div>

              <div class="bmol-rank-items-row">
                <div
                  v-for="item in getRankItemsToShow(group.items)"
                  :key="item.id"
                  class="bmol-item-card"
                >
                  <router-link
                    :to="item.media_type === 'tv' ? `/tv/${item.media.id}` : `/movies/${item.media.id}`"
                    class="bmol-poster-frame"
                  >
                    <img
                      v-if="item.media.poster_url"
                      :src="`${TMDB_IMG}${item.media.poster_url}`"
                      :alt="item.media.title"
                      loading="lazy"
                    />
                    <div v-else class="poster-fallback">
                      <Film :size="24" />
                    </div>
                  </router-link>

                  <div class="bmol-meta">
                    <router-link
                      :to="item.media_type === 'tv' ? `/tv/${item.media.id}` : `/movies/${item.media.id}`"
                      class="bmol-title"
                    >
                      {{ item.media.title }}
                    </router-link>
                    <span class="bmol-rating">★ {{ item.media.vote_average.toFixed(1) }}</span>

                    <!-- Actions (Only for Owner) -->
                    <div v-if="isOwner" class="bmol-actions">
                      <button
                        class="bmol-action-btn"
                        :title="$t('library.bmol.increaseRank')"
                        :disabled="item.rank <= 1"
                        @click="increaseRank(item)"
                      >
                        ▲
                      </button>
                      <button
                        class="bmol-action-btn"
                        :title="$t('library.bmol.decreaseRank')"
                        @click="decreaseRank(item)"
                      >
                        ▼
                      </button>
                      <button
                        class="bmol-action-btn bmol-action-btn--danger"
                        @click="triggerRemoveBmolItem(item)"
                      >
                        ✕
                      </button>
                    </div>
                  </div>
                </div>

                <!-- "+ X More" Card (If more than 3 items exist in this rank) -->
                <div
                  v-if="group.items.length > 3"
                  class="bmol-item-card bmol-more-card"
                  @click="openShowAllDetail(group.rank)"
                >
                  <div class="more-card-inner">
                    <span class="more-count">
                      {{ $t("library.bmol.moreItems", { count: group.items.length - 3 }) }}
                    </span>
                    <span class="more-sub">Click to view all</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Pagination Controls -->
            <div v-if="bmolTotalPages > 1" class="bmol-pagination">
              <button
                class="bmol-page-btn"
                :disabled="bmolPage <= 1"
                @click="setBmolPage(bmolPage - 1)"
              >
                ‹
              </button>
              <button
                v-for="p in bmolTotalPages"
                :key="p"
                class="bmol-page-btn"
                :class="{ 'bmol-page-btn--active': p === bmolPage }"
                @click="setBmolPage(p)"
              >
                {{ p }}
              </button>
              <button
                class="bmol-page-btn"
                :disabled="bmolPage >= bmolTotalPages"
                @click="setBmolPage(bmolPage + 1)"
              >
                ›
              </button>
            </div>
          </div>

          <div v-else class="empty-state-card">
            <Trophy :size="32" class="empty-icon" />
            <p>{{ $t("library.bmol.empty") }}</p>
          </div>
        </div>
      </section>

      <!-- TAB CONTENT: Lists (Watchlist, Watched, Likes, Reviews) -->
      <section v-else class="list-section fade-in-up">
        <!-- Advanced Filter Bar -->
        <div class="filter-bar">
          <!-- Title Search -->
          <div class="filter-input-wrapper search-wrapper">
            <Search :size="16" class="filter-icon" />
            <input
              v-model="filters.search"
              type="text"
              :placeholder="$t('library.filters.search')"
              class="filter-input"
            />
          </div>

          <!-- Media Type Select -->
          <div class="filter-select-wrapper">
            <Film :size="14" class="select-icon" />
            <select v-model="filters.mediaType" class="filter-select">
              <option value="all">{{ $t("library.filters.allMedia") }}</option>
              <option value="movie">{{ $t("library.filters.movies") }}</option>
              <option value="tv">{{ $t("library.filters.tvSeries") }}</option>
            </select>
          </div>

          <!-- Genre Filter -->
          <div class="filter-select-wrapper">
            <ArrowUpDown :size="14" class="select-icon" />
            <select v-model="filters.genre" class="filter-select">
              <option value="all">{{ $t("library.filters.allGenres") }}</option>
              <option
                v-for="genre in availableGenres"
                :key="genre"
                :value="genre"
              >{{ genre }}</option>
            </select>
          </div>

          <!-- Sort Order -->
          <div class="filter-select-wrapper">
            <ArrowUpDown :size="14" class="select-icon" />
            <select v-model="filters.sortBy" class="filter-select">
              <option value="newest">{{ $t("library.filters.sort.newest") }}</option>
              <option value="oldest">{{ $t("library.filters.sort.oldest") }}</option>
              <option v-if="activeTab === 'reviews'" value="rating_high">{{ $t("library.filters.sort.ratingHigh") }}</option>
              <option v-if="activeTab === 'reviews'" value="rating_low">{{ $t("library.filters.sort.ratingLow") }}</option>
              <option value="title_az">{{ $t("library.filters.sort.titleAZ") }}</option>
              <option value="title_za">{{ $t("library.filters.sort.titleZA") }}</option>
            </select>
          </div>
        </div>

        <!-- Render List Items -->
        <!-- Tab: Reviews -->
        <div v-if="activeTab === 'reviews'" class="reviews-list">
          <div
            v-for="review in filteredReviews"
            :key="review.id"
            class="review-library-card"
          >
            <router-link
              :to="review.media.media_type === 'tv' ? `/tv/${review.media.id}` : `/movies/${review.media.id}`"
              class="review-media-poster"
            >
              <img
                v-if="review.media.poster_url"
                :src="`${TMDB_IMG}${review.media.poster_url}`"
                :alt="review.media.title"
                loading="lazy"
              />
              <div v-else class="poster-fallback">
                <Film :size="18" />
              </div>
            </router-link>
            <div class="review-card-content">
              <div class="review-card-header">
                <router-link
                  :to="review.media.media_type === 'tv' ? `/tv/${review.media.id}` : `/movies/${review.media.id}`"
                  class="review-media-title"
                >
                  {{ review.media.title }}
                </router-link>
                <div class="review-card-stars">
                  <span class="rating-num">{{ review.rating.toFixed(1) }}</span>
                  <span class="stars-icon">★</span>
                </div>
              </div>
              <p class="review-card-body">{{ review.body }}</p>
              <div class="review-card-footer">
                <span class="review-date">{{ new Date(review.created_at).toLocaleDateString() }}</span>
                <div class="review-meta-actions">
                  <span class="action-count"><Heart :size="12" class="liked-heart" /> {{ review.like_count }}</span>
                </div>
              </div>
            </div>
          </div>
          <div v-if="filteredReviews.length === 0" class="empty-state-card">
            <Star :size="32" class="empty-icon" />
            <p>{{ $t("library.empty_reviews") }}</p>
          </div>
        </div>

        <!-- Tab: Library Items (Watchlist, Watched, Likes) -->
        <div v-else class="media-grid">
          <div
            v-for="item in filteredLibraryItems"
            :key="item.id"
            class="media-library-card"
          >
            <router-link
              :to="item.media.media_type === 'tv' ? `/tv/${item.media.id}` : `/movies/${item.media.id}`"
              class="media-poster-frame"
            >
              <img
                v-if="item.media.poster_url"
                :src="`${TMDB_IMG}${item.media.poster_url}`"
                :alt="item.media.title"
                loading="lazy"
              />
              <div v-else class="poster-fallback">
                <Film :size="24" />
              </div>
              <span class="media-type-chip">{{ item.media.media_type === 'tv' ? 'TV' : 'Movie' }}</span>
            </router-link>
            <div class="media-card-meta">
              <router-link
                :to="item.media.media_type === 'tv' ? `/tv/${item.media.id}` : `/movies/${item.media.id}`"
                class="media-title"
              >
                {{ item.media.title }}
              </router-link>
              <div class="media-details-row">
                <span class="media-rating-badge">★ {{ item.media.vote_average.toFixed(1) }}</span>
                <span class="save-date">{{ new Date(item.created_at).toLocaleDateString() }}</span>
              </div>
              <!-- User tags if present -->
              <div v-if="item.tags && item.tags.length > 0" class="tags-wrapper">
                <span v-for="t in item.tags" :key="t" class="tag-badge">{{ t }}</span>
              </div>
            </div>
          </div>
          <div v-if="filteredLibraryItems.length === 0" class="empty-state-card">
            <Bookmark :size="32" class="empty-icon" />
            <p>{{ $t("library.empty") }}</p>
          </div>
        </div>
      </section>
    </div>

    <!-- Confirm Modal for Deletion -->
    <ConfirmModal
      v-model="showDeleteConfirm"
      list-type="bmol_delete"
      :item-name="deleteItemName"
      @confirm="confirmDelete"
      @cancel="showDeleteConfirm = false"
    />

    <!-- SPOTLIGHT QUICK ADD MODAL (Teleported to body) -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="spotlightActive" class="spotlight-backdrop" @click.self="closeSpotlight">
          <div class="spotlight-modal">
            
            <!-- SELECTING ZONE (Top) -->
            <div v-if="selectedSpotlightItems.length > 0" class="spotlight-selecting-zone">
              <div
                v-for="item in selectedSpotlightItems"
                :key="item.id"
                class="selecting-chip"
              >
                <img
                  v-if="item.poster_path"
                  :src="`${TMDB_IMG}${item.poster_path}`"
                  :alt="getMediaTitle(item)"
                  class="selecting-chip-poster"
                />
                <span class="selecting-chip-title">{{ getMediaTitle(item) }}</span>
                <button class="btn-remove-selecting" @click="unpinSpotlightItem(item.id)">✕</button>
              </div>
            </div>

            <!-- SPOTLIGHT BAR (Middle) -->
            <div class="spotlight-bar">
              <div class="spotlight-search-wrapper">
                <Search :size="18" class="spotlight-search-icon" />
                <input
                  ref="spotlightInput"
                  v-model="spotlightQuery"
                  type="text"
                  :placeholder="$t('library.bmol.searchPlaceholder')"
                  class="spotlight-search-input"
                />
              </div>

              <!-- Select Multiple Mode Checkbox/Toggle -->
              <div class="spotlight-toggle-wrapper">
                <label class="switch-container">
                  <input type="checkbox" v-model="selectMultipleMode" />
                  <span class="switch-slider" />
                </label>
                <span class="switch-label">{{ $t("library.bmol.selectMultiple") }}</span>
              </div>
            </div>

            <!-- SEARCH RESULTS (Bottom) -->
            <div class="spotlight-results">
              <div v-if="spotlightSearchLoading" class="spotlight-loading text-center py-4">
                <div class="loading-ring small-ring mx-auto" />
              </div>
              <div v-else-if="spotlightResults.length > 0" class="spotlight-results-list">
                <div
                  v-for="res in spotlightResults"
                  :key="res.id"
                  class="spotlight-result-row"
                  :class="{ 'spotlight-result-row--selected': isSpotlightItemSelected(res.id) }"
                  @click="handleResultClick(res)"
                >
                  <img
                    v-if="res.poster_path"
                    :src="`${TMDB_IMG}${res.poster_path}`"
                    :alt="getMediaTitle(res)"
                    class="result-row-poster"
                  />
                  <div class="result-row-poster-fallback" v-else>
                    <Film :size="16" />
                  </div>
                  <div class="result-row-info">
                    <p class="result-row-title">{{ getMediaTitle(res) }}</p>
                    <p class="result-row-date">{{ getMediaDate(res) || '—' }}</p>
                  </div>
                  <div class="result-row-indicator">
                    <span v-if="isSpotlightItemSelected(res.id)" class="indicator-check">✓</span>
                    <span v-else class="indicator-plus">+</span>
                  </div>
                </div>
              </div>
              <div v-else-if="spotlightQuery.trim() !== ''" class="spotlight-empty-results">
                <p>{{ $t("library.bmol.errorEmptySearch") }}</p>
              </div>
            </div>

            <!-- MODAL ACTION FOOTER -->
            <div class="spotlight-footer">
              <button class="btn-spotlight-cancel" @click="closeSpotlight">Close</button>
              <button
                v-if="selectMultipleMode"
                class="btn-spotlight-save"
                :disabled="selectedSpotlightItems.length === 0"
                @click="saveSpotlightItems"
              >
                {{ $t("library.bmol.addSelected", { count: selectedSpotlightItems.length }) }}
              </button>
            </div>

          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useRoute } from "vue-router"
import { useAuthStore } from "@/stores/auth"
import { useI18n } from "vue-i18n"
import { libraryApi, reviewApi, userApi, bmolApi, movieApi } from "@/api/api"
import type { UserProfile } from "@/types/user"
import ConfirmModal from "@/components/profile/components/ConfirmModal.vue"
import type { LibraryItemResponse, ReviewResponse, BMOLItemResponse, Movie, TVSeries } from "@/types"
import {
  Search,
  Bookmark,
  Eye,
  Heart,
  Star,
  ArrowUpDown,
  Flame,
  Trophy,
  ChevronLeft,
  User as UserIcon,
  Film,
  TrendingUp,
  TrendingDown,
  Minus
} from "lucide-vue-next"

const route = useRoute()
const auth = useAuthStore()
const { t } = useI18n()

const TMDB_IMG = "https://image.tmdb.org/t/p/w342"

const userId = computed(() => {
  const raw = route.params.userId
  return Number(Array.isArray(raw) ? raw[0] : raw)
})

const loading = ref(true)
const activeTab = ref<"dashboard" | "bmol" | "watchlist" | "watched" | "likes" | "reviews">("dashboard")
const bmolSubTab = ref<"movie" | "tv">("movie")

const profileUser = ref<UserProfile | null>(null)
const watchlistItems = ref<LibraryItemResponse[]>([])
const watchedItems = ref<LibraryItemResponse[]>([])
const likedItems = ref<LibraryItemResponse[]>([])
const userReviews = ref<ReviewResponse[]>([])
const bmolItems = ref<BMOLItemResponse[]>([])

const isOwner = computed(() => auth.user?.id === userId.value)

// Filter States
const filters = ref({
  search: "",
  mediaType: "all",
  genre: "all",
  sortBy: "newest"
})

// Load everything on mount
const loadData = async () => {
  try {
    loading.value = true
    const id = userId.value

    // Fetch profile info
    const profileRes = await userApi.getProfile(id)
    profileUser.value = profileRes.data.user

    // Fetch lists in parallel
    const [watchlistRes, watchedRes, likedRes, reviewsRes, bmolRes] = await Promise.all([
      libraryApi.getVisibleUserLibrary(id, { list_type: "watchlist" }),
      libraryApi.getVisibleUserLibrary(id, { list_type: "watched" }),
      libraryApi.getVisibleUserLibrary(id, { list_type: "likes" }),
      reviewApi.getUserReviews(id),
      bmolApi.getUserBMOL(id)
    ])

    watchlistItems.value = watchlistRes.data.items
    watchedItems.value = watchedRes.data.items
    likedItems.value = likedRes.data.items
    userReviews.value = reviewsRes.data.reviews
    bmolItems.value = bmolRes.data.items
  } catch (err) {
    console.error("Failed to load library data:", err)
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
watch(userId, loadData)

// BMOL Spotlight & Detail Logic
const selectedRankDetail = ref<number | null>(null)
const rankFilterQuery = ref("")

const bmolFilteredItems = computed(() => {
  if (selectedRankDetail.value === null) return []
  const rankItems = bmolItems.value.filter(
    item => item.rank === selectedRankDetail.value &&
            item.media_type === (bmolSubTab.value === "movie" ? "movie" : "tv")
  )
  const query = rankFilterQuery.value.trim().toLowerCase()
  if (!query) return rankItems
  return rankItems.filter(item =>
    item.media.title.toLowerCase().includes(query)
  )
})

function openShowAllDetail(rank: number) {
  selectedRankDetail.value = rank
  rankFilterQuery.value = ""
}

function goBackToList() {
  selectedRankDetail.value = null
  rankFilterQuery.value = ""
}

const showDeleteConfirm = ref(false)
const deleteItemId = ref<number | null>(null)
const deleteItemName = ref("")

function triggerRemoveBmolItem(item: BMOLItemResponse) {
  deleteItemId.value = item.id
  deleteItemName.value = item.media.title
  showDeleteConfirm.value = true
}

function getRankItemsToShow(items: BMOLItemResponse[]) {
  return items.slice(0, 3)
}

// --- Pagination ---
const BMOL_PER_PAGE = 5
const bmolPage = ref(1)

const bmolTotalPages = computed(() =>
  Math.max(1, Math.ceil(groupedBmolItems.value.length / BMOL_PER_PAGE))
)

const paginatedBmolGroups = computed(() => {
  const start = (bmolPage.value - 1) * BMOL_PER_PAGE
  return groupedBmolItems.value.slice(start, start + BMOL_PER_PAGE)
})

function setBmolPage(page: number) {
  if (page >= 1 && page <= bmolTotalPages.value) {
    bmolPage.value = page
  }
}

// --- Add Form (Search bar to add media to any rank) ---
const bmolAddQuery = ref("")
const bmolAddResults = ref<Array<Movie | TVSeries>>([])
const bmolAddLoading = ref(false)
const bmolAddSelected = ref<(Movie | TVSeries) | null>(null)
const bmolAddRank = ref(1)
let bmolAddSearchTimeout: ReturnType<typeof setTimeout> | null = null

watch(bmolAddQuery, (newVal) => {
  if (bmolAddSearchTimeout) clearTimeout(bmolAddSearchTimeout)
  if (!newVal.trim()) {
    bmolAddResults.value = []
    return
  }
  bmolAddLoading.value = true
  bmolAddSearchTimeout = setTimeout(async () => {
    try {
      if (bmolSubTab.value === "movie") {
        const res = await movieApi.search(newVal)
        bmolAddResults.value = res.data.results.slice(0, 6)
      } else {
        const res = await movieApi.searchSeries(newVal)
        bmolAddResults.value = res.data.results.slice(0, 6)
      }
    } catch (err) {
      console.error("BMOL add search failed:", err)
    } finally {
      bmolAddLoading.value = false
    }
  }, 350)
})

function selectBmolAddResult(media: Movie | TVSeries) {
  bmolAddSelected.value = media
  bmolAddQuery.value = ""
  bmolAddResults.value = []
}

function clearBmolAddSelection() {
  bmolAddSelected.value = null
}

async function submitBmolAdd() {
  if (!bmolAddSelected.value) return
  const media = bmolAddSelected.value
  const rank = bmolAddRank.value

  // Optimistic UI
  const tempItem: BMOLItemResponse = {
    id: Date.now() + Math.random(),
    rank,
    media_id: media.id,
    media_type: bmolSubTab.value,
    media: {
      id: media.id,
      title: getMediaTitle(media),
      poster_url: media.poster_path || "",
      vote_average: media.vote_average || 0
    }
  } as BMOLItemResponse
  bmolItems.value.push(tempItem)
  clearBmolAddSelection()

  try {
    await bmolApi.addItem({
      media_id: media.id,
      media_type: bmolSubTab.value,
      rank
    })
    // Refresh from server
    const res = await bmolApi.getUserBMOL(userId.value)
    bmolItems.value = res.data.items
  } catch (err: unknown) {
    const errorWithResponse = err as { response?: { status: number } }
    if (errorWithResponse.response?.status === 409) {
      bmolItems.value = bmolItems.value.filter(i => i.media_id !== media.id || i.id !== tempItem.id)
    } else {
      bmolItems.value = bmolItems.value.filter(i => i.id !== tempItem.id)
      console.error("Failed to add BMOL item:", err)
    }
  }
}

// Spotlight State
const spotlightActive = ref(false)
const spotlightRank = ref<number | null>(null)
const spotlightQuery = ref("")
const spotlightResults = ref<Array<Movie | TVSeries>>([])
const spotlightSearchLoading = ref(false)
const selectedSpotlightItems = ref<Array<Movie | TVSeries>>([])
const selectMultipleMode = ref(false)
const spotlightInput = ref<HTMLInputElement | null>(null)

function openSpotlight(rank: number) {
  spotlightRank.value = rank
  spotlightQuery.value = ""
  spotlightResults.value = []
  selectedSpotlightItems.value = []
  selectMultipleMode.value = false
  spotlightActive.value = true
  setTimeout(() => {
    spotlightInput.value?.focus()
  }, 100)
}

function closeSpotlight() {
  spotlightActive.value = false
  spotlightRank.value = null
  spotlightQuery.value = ""
  spotlightResults.value = []
  selectedSpotlightItems.value = []
}

function isSpotlightItemSelected(id: number) {
  return selectedSpotlightItems.value.some(item => item.id === id)
}

function unpinSpotlightItem(id: number) {
  selectedSpotlightItems.value = selectedSpotlightItems.value.filter(item => item.id !== id)
}

function handleResultClick(media: Movie | TVSeries) {
  if (selectMultipleMode.value) {
    if (isSpotlightItemSelected(media.id)) {
      unpinSpotlightItem(media.id)
    } else {
      selectedSpotlightItems.value.push(media)
    }
  } else {
    selectedSpotlightItems.value = [media]
    saveSpotlightItems()
  }
}

async function saveSpotlightItems() {
  if (selectedSpotlightItems.value.length === 0 || spotlightRank.value === null) return
  const rank = spotlightRank.value
  const itemsToSave = [...selectedSpotlightItems.value]
  closeSpotlight()

  // Optimistic UI update
  const newBmolItems = itemsToSave.map(media => {
    return {
      id: Date.now() + Math.random(),
      rank: rank,
      media_id: media.id,
      media_type: bmolSubTab.value,
      media: {
        id: media.id,
        title: getMediaTitle(media),
        poster_url: media.poster_path || "",
        vote_average: media.vote_average || 0
      }
    } as BMOLItemResponse
  })
  bmolItems.value.push(...newBmolItems)

  for (const media of itemsToSave) {
    try {
      await bmolApi.addItem({
        media_id: media.id,
        media_type: bmolSubTab.value,
        rank: rank
      })
    } catch (err: unknown) {
      const errorWithResponse = err as { response?: { status: number } }
      if (errorWithResponse.response?.status === 409) {
        bmolItems.value = bmolItems.value.filter(i => i.media_id !== media.id)
      } else {
        bmolItems.value = bmolItems.value.filter(i => i.media_id !== media.id)
        console.error("Failed to add spotlight item:", err)
      }
    }
  }

  // Refresh
  try {
    const res = await bmolApi.getUserBMOL(userId.value)
    bmolItems.value = res.data.items
  } catch (err) {
    console.error("Refresh failed:", err)
  }
}

let spotlightSearchTimeout: ReturnType<typeof setTimeout> | null = null
watch(spotlightQuery, (newVal) => {
  if (spotlightSearchTimeout) clearTimeout(spotlightSearchTimeout)
  if (!newVal.trim()) {
    spotlightResults.value = []
    return
  }
  spotlightSearchLoading.value = true
  spotlightSearchTimeout = setTimeout(async () => {
    try {
      if (bmolSubTab.value === "movie") {
        const res = await movieApi.search(newVal)
        spotlightResults.value = res.data.results.slice(0, 5)
      } else {
        const res = await movieApi.searchSeries(newVal)
        spotlightResults.value = res.data.results.slice(0, 5)
      }
    } catch (err) {
      console.error("Spotlight search failed:", err)
    } finally {
      spotlightSearchLoading.value = false
    }
  }, 350)
})

watch(spotlightActive, (newVal) => {
  if (newVal) {
    document.body.style.overflow = "hidden"
  } else {
    const activeBackdrops = document.querySelectorAll(".modal-backdrop, .spotlight-backdrop")
    if (activeBackdrops.length <= 1) {
      document.body.style.overflow = ""
    }
  }
})

watch(bmolSubTab, () => {
  selectedRankDetail.value = null
  bmolPage.value = 1
  bmolAddQuery.value = ""
  bmolAddResults.value = []
  bmolAddSelected.value = null
  closeSpotlight()
})

function getMediaTitle(media: Movie | TVSeries): string {
  if ("title" in media) {
    return media.title
  }
  return media.name
}

function getMediaDate(media: Movie | TVSeries): string {
  if ("release_date" in media) {
    return media.release_date
  }
  return media.first_air_date
}

async function increaseRank(item: BMOLItemResponse) {
  if (item.rank <= 1) return
  const oldRank = item.rank
  const newRank = item.rank - 1
  const found = bmolItems.value.find(i => i.id === item.id)
  if (found) found.rank = newRank

  try {
    await bmolApi.updateItem(item.id, { rank: newRank })
  } catch (err) {
    const foundRollback = bmolItems.value.find(i => i.id === item.id)
    if (foundRollback) foundRollback.rank = oldRank
    console.error("Failed to increase rank:", err)
  }
}

async function decreaseRank(item: BMOLItemResponse) {
  const oldRank = item.rank
  const newRank = item.rank + 1
  const found = bmolItems.value.find(i => i.id === item.id)
  if (found) found.rank = newRank

  try {
    await bmolApi.updateItem(item.id, { rank: newRank })
  } catch (err) {
    const foundRollback = bmolItems.value.find(i => i.id === item.id)
    if (foundRollback) foundRollback.rank = oldRank
    console.error("Failed to decrease rank:", err)
  }
}

async function confirmDelete() {
  if (deleteItemId.value === null) return
  const itemId = deleteItemId.value
  showDeleteConfirm.value = false

  const backup = [...bmolItems.value]
  bmolItems.value = bmolItems.value.filter(i => i.id !== itemId)

  try {
    await bmolApi.removeItem(itemId)
  } catch (err) {
    bmolItems.value = backup
    console.error("Failed to remove item:", err)
  } finally {
    deleteItemId.value = null
    deleteItemName.value = ""
  }
}

const groupedBmolItems = computed(() => {
  const filtered = bmolItems.value.filter(
    item => item.media_type === (bmolSubTab.value === "movie" ? "movie" : "tv")
  )
  const groups: Record<number, BMOLItemResponse[]> = {}
  filtered.forEach(item => {
    if (!groups[item.rank]) {
      groups[item.rank] = []
    }
    groups[item.rank]?.push(item)
  })

  return Object.entries(groups)
    .map(([rankStr, items]) => ({
      rank: Number(rankStr),
      items
    }))
    .sort((a, b) => a.rank - b.rank)
})

// Dashboard Computations
const averageRating = computed(() => {
  if (userReviews.value.length === 0) return 0
  const sum = userReviews.value.reduce((acc, r) => acc + r.rating, 0)
  return Number((sum / userReviews.value.length).toFixed(1))
})

const topGenresList = computed(() => {
  const counts: Record<string, number> = {}
  const allMediaItems = [...watchedItems.value, ...likedItems.value]

  allMediaItems.forEach(item => {
    if (item.media && item.media.genres) {
      item.media.genres.forEach(g => {
        counts[g.name] = (counts[g.name] || 0) + 1
      })
    }
  })

  return Object.entries(counts)
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 5)
})

const topGenre = computed(() => {
  return topGenresList.value[0]?.name || null
})

const reviewsStats = computed(() => {
  const now = new Date()
  const currentYear = now.getFullYear()
  const currentMonth = now.getMonth()

  let thisMonthCount = 0
  let lastMonthCount = 0

  userReviews.value.forEach(r => {
    const rDate = new Date(r.created_at)
    const rYear = rDate.getFullYear()
    const rMonth = rDate.getMonth()

    if (rYear === currentYear && rMonth === currentMonth) {
      thisMonthCount++
    } else if (
      (currentMonth === 0 && rYear === currentYear - 1 && rMonth === 11) ||
      (currentMonth > 0 && rYear === currentYear && rMonth === currentMonth - 1)
    ) {
      lastMonthCount++
    }
  })

  const diff = thisMonthCount - lastMonthCount
  return { thisMonth: thisMonthCount, lastMonth: lastMonthCount, diff }
})

const trendIcon = computed(() => {
  if (reviewsStats.value.diff > 0) return TrendingUp
  if (reviewsStats.value.diff < 0) return TrendingDown
  return Minus
})

const trendClass = computed(() => {
  if (reviewsStats.value.diff > 0) return "trend-green"
  if (reviewsStats.value.diff < 0) return "trend-red"
  return "trend-neutral"
})

const ratingDistribution = computed(() => {
  const distribution = Array(10).fill(0) // 0.5 to 5.0 in steps of 0.5
  userReviews.value.forEach(r => {
    const index = Math.round(r.rating * 2) - 1
    if (index >= 0 && index < 10) {
      distribution[index]++
    }
  })
  const maxCount = Math.max(...distribution, 1)
  return distribution
    .map((count, i) => ({
      rating: (i + 1) / 2,
      count,
      percentage: (count / maxCount) * 100
    }))
    .reverse()
})

// Dynamic Genres for Dropdown Filter
const availableGenres = computed(() => {
  const genresMap = new Set<string>()
  const items =
    activeTab.value === "reviews"
      ? userReviews.value.map(r => r.media)
      : activeTab.value === "watchlist"
      ? watchlistItems.value.map(i => i.media)
      : activeTab.value === "watched"
      ? watchedItems.value.map(i => i.media)
      : likedItems.value.map(i => i.media)

  items.forEach(media => {
    if (media && media.genres) {
      media.genres.forEach(g => genresMap.add(g.name))
    }
  })

  return Array.from(genresMap).sort()
})

// Filter & Sort Logic for Library Items
const filteredLibraryItems = computed(() => {
  let list =
    activeTab.value === "watchlist"
      ? watchlistItems.value
      : activeTab.value === "watched"
      ? watchedItems.value
      : likedItems.value

  // 1. Search Query
  if (filters.value.search.trim()) {
    const query = filters.value.search.toLowerCase()
    list = list.filter(item => item.media.title.toLowerCase().includes(query))
  }

  // 2. Media Type
  if (filters.value.mediaType !== "all") {
    list = list.filter(item => item.media.media_type === filters.value.mediaType)
  }

  // 3. Genre
  if (filters.value.genre !== "all") {
    list = list.filter(
      item => item.media.genres && item.media.genres.some(g => g.name === filters.value.genre)
    )
  }

  // 4. Sort
  const sorted = [...list]
  if (filters.value.sortBy === "newest") {
    sorted.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  } else if (filters.value.sortBy === "oldest") {
    sorted.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
  } else if (filters.value.sortBy === "title_az") {
    sorted.sort((a, b) => a.media.title.localeCompare(b.media.title))
  } else if (filters.value.sortBy === "title_za") {
    sorted.sort((a, b) => b.media.title.localeCompare(a.media.title))
  }

  return sorted
})

// Filter & Sort Logic for Reviews
const filteredReviews = computed(() => {
  let list = userReviews.value

  // 1. Search Query
  if (filters.value.search.trim()) {
    const query = filters.value.search.toLowerCase()
    list = list.filter(r => r.media.title.toLowerCase().includes(query))
  }

  // 2. Media Type
  if (filters.value.mediaType !== "all") {
    list = list.filter(r => r.media.media_type === filters.value.mediaType)
  }

  // 3. Genre
  if (filters.value.genre !== "all") {
    list = list.filter(
      r => r.media.genres && r.media.genres.some(g => g.name === filters.value.genre)
    )
  }

  // 4. Sort
  const sorted = [...list]
  if (filters.value.sortBy === "newest") {
    sorted.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  } else if (filters.value.sortBy === "oldest") {
    sorted.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
  } else if (filters.value.sortBy === "rating_high") {
    sorted.sort((a, b) => b.rating - a.rating)
  } else if (filters.value.sortBy === "rating_low") {
    sorted.sort((a, b) => a.rating - b.rating)
  } else if (filters.value.sortBy === "title_az") {
    sorted.sort((a, b) => a.media.title.localeCompare(b.media.title))
  } else if (filters.value.sortBy === "title_za") {
    sorted.sort((a, b) => b.media.title.localeCompare(a.media.title))
  }

  return sorted
})

// Reset Filters on Tab Change
watch(activeTab, () => {
  filters.value.search = ""
  filters.value.mediaType = "all"
  filters.value.genre = "all"
  filters.value.sortBy = "newest"
})
</script>

<style scoped>
.library-root {
  font-family: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background-color: #0c0c0e;
  color: #f3f4f6;
  min-height: 100vh;
  padding: 2.5rem 2rem;
  position: relative;
  overflow: hidden;
}

.grain {
  position: fixed;
  inset: 0;
  background-image: radial-gradient(rgba(255, 255, 255, 0.015) 1px, transparent 0);
  background-size: 24px 24px;
  pointer-events: none;
  z-index: 1;
}

/* Loading styling */
.library-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 70vh;
  gap: 1rem;
}

.loading-ring {
  width: 48px;
  height: 48px;
  border: 3px solid rgba(225, 37, 27, 0.1);
  border-top-color: #e1251b;
  border-radius: 50%;
  animation: spin 1s infinite linear;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-text {
  font-size: 0.95rem;
  color: #9ca3af;
}

.library-layout {
  position: relative;
  z-index: 2;
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

/* Back Link */
.back-link {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  color: #9ca3af;
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 500;
  transition: color 0.2s;
  align-self: flex-start;
}

.back-link:hover {
  color: #e1251b;
}

/* Summary Card Header */
.profile-summary-card {
  background: rgba(25, 25, 28, 0.65);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

@media (min-width: 768px) {
  .profile-summary-card {
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
  }
}

.user-info-section {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.avatar-wrapper {
  position: relative;
}

.avatar-ring {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e1251b, #8a1612);
  padding: 2px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.avatar-fallback {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: #27272a;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
}

.level-chip {
  position: absolute;
  bottom: -4px;
  right: -4px;
  background: #e1251b;
  color: #ffffff;
  font-size: 0.675rem;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 8px;
  border: 2px solid #161618;
}

.user-meta {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-size: 1.5rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
  letter-spacing: -0.01em;
}

.user-tag {
  font-size: 0.875rem;
  color: #9ca3af;
  margin: 0.125rem 0 0.5rem 0;
}

.user-bio {
  font-size: 0.875rem;
  color: #d1d5db;
  margin: 0;
  max-width: 500px;
  line-height: 1.5;
}

.user-counts-section {
  display: flex;
  gap: 2rem;
}

.count-box {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.count-val {
  font-size: 1.5rem;
  font-weight: 800;
  color: #ffffff;
}

.count-lbl {
  font-size: 0.75rem;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-top: 0.25rem;
}

/* Tabs */
.library-tabs {
  display: flex;
  gap: 0.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  padding-bottom: 0.5rem;
  overflow-x: auto;
}

.tab-btn {
  background: transparent;
  border: none;
  color: #9ca3af;
  padding: 0.75rem 1.25rem;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border-radius: 10px;
  transition: background-color 0.2s, color 0.2s;
  white-space: nowrap;
}

.tab-btn:hover {
  background: rgba(255, 255, 255, 0.03);
  color: #ffffff;
}

.tab-btn--active {
  background: rgba(225, 37, 27, 0.1);
  color: #e1251b;
}

.tab-count {
  font-size: 0.75rem;
  background: rgba(255, 255, 255, 0.08);
  color: #d1d5db;
  padding: 2px 6px;
  border-radius: 6px;
  margin-left: 0.25rem;
}

.tab-btn--active .tab-count {
  background: rgba(225, 37, 27, 0.2);
  color: #e1251b;
}

/* Dashboard styling */
.dashboard-section {
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.5rem;
}

@media (min-width: 640px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

.stat-card {
  background: rgba(20, 20, 22, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 16px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1.25rem;
  backdrop-filter: blur(8px);
}

.stat-icon-box {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.watched-theme { color: #3b82f6; background: rgba(59, 130, 246, 0.1); }
.review-theme { color: #e1251b; background: rgba(225, 37, 27, 0.1); }
.rating-theme { color: #ffb800; background: rgba(255, 184, 0, 0.1); }
.genre-theme { color: #a855f7; background: rgba(168, 85, 247, 0.1); }

.stat-data h3 {
  font-size: 0.8125rem;
  color: #9ca3af;
  margin: 0 0 0.25rem 0;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 800;
  color: #ffffff;
  margin: 0;
}

.genre-text {
  font-size: 1.125rem;
  font-weight: 700;
}

.trend-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 2px 6px;
  border-radius: 6px;
  margin-top: 0.375rem;
  font-size: 0.725rem;
  font-weight: 600;
}

.trend-green { background: rgba(16, 185, 129, 0.1); color: #10b981; }
.trend-red { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
.trend-neutral { background: rgba(255, 255, 255, 0.05); color: #9ca3af; }

.rating-display {
  display: flex;
  align-items: baseline;
  gap: 0.125rem;
}

.max-val {
  font-size: 0.8125rem;
  color: #9ca3af;
}

.stars-row {
  display: flex;
  gap: 0.125rem;
  margin-top: 0.25rem;
}

.star-dot {
  font-size: 0.875rem;
  color: #3f3f46;
}

.star-dot--active {
  color: #ffb800;
}

/* Charts Section */
.charts-row {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
}

@media (min-width: 1024px) {
  .charts-row {
    grid-template-columns: 1fr 1fr;
  }
}

.chart-card {
  background: rgba(20, 20, 22, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 20px;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.chart-title {
  font-size: 1.125rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
  border-left: 3px solid #e1251b;
  padding-left: 0.75rem;
}

/* Rating Distribution CSS Chart */
.distribution-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.dist-row {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.dist-label {
  font-size: 0.8125rem;
  color: #9ca3af;
  width: 36px;
  text-align: right;
  font-weight: 600;
}

.dist-bar-wrapper {
  flex: 1;
  height: 8px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 4px;
  overflow: hidden;
}

.dist-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #8a1612, #e1251b);
  border-radius: 4px;
}

.dist-count {
  font-size: 0.8125rem;
  color: #ffffff;
  width: 24px;
  font-weight: 700;
}

/* Top Genres List ranking */
.genres-ranking {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.genre-rank-row {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.rank-badge {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.8125rem;
  background: rgba(255, 255, 255, 0.05);
  color: #9ca3af;
}

.rank-1 { background: rgba(255, 184, 0, 0.15); color: #ffb800; }
.rank-2 { background: rgba(156, 163, 175, 0.15); color: #d1d5db; }
.rank-3 { background: rgba(180, 83, 9, 0.15); color: #b45309; }

.genre-rank-name {
  flex: 1;
  font-size: 0.875rem;
  font-weight: 600;
  color: #f3f4f6;
}

.genre-rank-count {
  font-size: 0.8125rem;
  color: #9ca3af;
}

.genres-empty {
  text-align: center;
  color: #52525b;
  font-size: 0.875rem;
  padding: 2rem 0;
}

/* BMOL section styling */
.bmol-section {
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

/* ==============================
   BMOL Add Form
   ============================== */
.bmol-add-card {
  background: rgba(20, 20, 22, 0.55);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 20px;
  backdrop-filter: blur(12px);
  margin-bottom: 1.5rem;
}

.bmol-add-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: #e4e4e7;
  margin: 0 0 14px 0;
  padding-left: 10px;
  border-left: 2px solid #e1251b;
}

.bmol-add-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Search wrapper */
.bmol-add-search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  padding: 0 14px;
  height: 42px;
  transition: border-color 0.2s, background-color 0.2s;
}

.bmol-add-search-wrapper:focus-within {
  border-color: rgba(225, 37, 27, 0.35);
  background: rgba(255, 255, 255, 0.04);
}

.bmol-add-search-icon {
  color: #52525b;
  flex-shrink: 0;
  transition: color 0.2s;
}

.bmol-add-search-wrapper:focus-within .bmol-add-search-icon {
  color: #e1251b;
}

.bmol-add-search-input {
  background: transparent;
  border: none;
  outline: none;
  color: #e4e4e7;
  font-size: 0.85rem;
  font-weight: 500;
  width: 100%;
}

.bmol-add-search-input::placeholder {
  color: #52525b;
}

.bmol-add-spinner {
  flex-shrink: 0;
}

/* Dropdown */
.bmol-add-dropdown {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  background: #18181b;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  z-index: 20;
  max-height: 280px;
  overflow-y: auto;
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.5);
}

.bmol-add-dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  cursor: pointer;
  transition: background-color 0.15s;
  border-bottom: 1px solid rgba(255, 255, 255, 0.02);
}

.bmol-add-dropdown-item:last-child {
  border-bottom: none;
}

.bmol-add-dropdown-item:hover {
  background: rgba(255, 255, 255, 0.04);
}

.bmol-add-dropdown-poster {
  width: 30px;
  height: 45px;
  border-radius: 4px;
  object-fit: cover;
  flex-shrink: 0;
}

.bmol-add-dropdown-poster-fallback {
  width: 30px;
  height: 45px;
  border-radius: 4px;
  background: #27272a;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #52525b;
  flex-shrink: 0;
}

.bmol-add-dropdown-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.bmol-add-dropdown-title {
  font-size: 0.82rem;
  font-weight: 650;
  color: #ffffff;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bmol-add-dropdown-date {
  font-size: 0.7rem;
  color: #71717a;
  margin: 3px 0 0;
}

/* Selected preview */
.bmol-add-selected {
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  padding: 14px;
}

.bmol-add-selected-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.bmol-add-selected-poster {
  width: 40px;
  height: 60px;
  border-radius: 6px;
  object-fit: cover;
  flex-shrink: 0;
}

.bmol-add-selected-detail {
  flex: 1;
  min-width: 0;
}

.bmol-add-selected-title {
  font-size: 0.88rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bmol-add-selected-date {
  font-size: 0.72rem;
  color: #71717a;
  margin: 4px 0 0;
}

.bmol-add-clear-btn {
  background: rgba(255, 255, 255, 0.05);
  border: none;
  color: #71717a;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 0.7rem;
  flex-shrink: 0;
  transition: all 0.2s;
}

.bmol-add-clear-btn:hover {
  background: rgba(239, 68, 68, 0.12);
  color: #ef4444;
}

/* Rank selector row */
.bmol-add-rank-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.bmol-add-rank-label {
  font-size: 0.78rem;
  font-weight: 600;
  color: #a1a1aa;
  white-space: nowrap;
}

.bmol-add-rank-input {
  width: 60px;
  background: #1c1c1e;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  color: #ffffff;
  font-size: 0.85rem;
  font-weight: 600;
  padding: 6px 10px;
  text-align: center;
  outline: none;
  transition: border-color 0.2s;
}

.bmol-add-rank-input:focus {
  border-color: rgba(225, 37, 27, 0.4);
}

/* Chrome number input arrows */
.bmol-add-rank-input::-webkit-inner-spin-button,
.bmol-add-rank-input::-webkit-outer-spin-button {
  opacity: 1;
}

.bmol-add-submit-btn {
  margin-left: auto;
  background: #e1251b;
  color: #ffffff;
  border: none;
  border-radius: 8px;
  font-size: 0.78rem;
  font-weight: 650;
  padding: 8px 18px;
  cursor: pointer;
  transition: background-color 0.2s, transform 0.1s;
  white-space: nowrap;
}

.bmol-add-submit-btn:hover {
  background: #b81d15;
}

.bmol-add-submit-btn:active {
  transform: scale(0.97);
}

/* ==============================
   BMOL Pagination
   ============================== */
.bmol-pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 1.5rem;
  padding: 10px 0;
}

.bmol-page-btn {
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(255, 255, 255, 0.02);
  color: #a1a1aa;
  font-size: 0.82rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.bmol-page-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.06);
  color: #ffffff;
  border-color: rgba(255, 255, 255, 0.12);
}

.bmol-page-btn--active {
  background: rgba(225, 37, 27, 0.12) !important;
  color: #e1251b !important;
  border-color: rgba(225, 37, 27, 0.3) !important;
}

.bmol-page-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

/* Quick Add Button */
.btn-quick-add {
  background: rgba(225, 37, 27, 0.08);
  border: 1px solid rgba(225, 37, 27, 0.2);
  border-radius: 8px;
  color: #e1251b;
  font-size: 0.75rem;
  font-weight: 650;
  padding: 4px 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-quick-add:hover {
  background: #e1251b;
  color: #ffffff;
  border-color: #e1251b;
}

/* + X More Card styling */
.bmol-more-card {
  cursor: pointer;
  background: rgba(255, 255, 255, 0.01) !important;
  border: 2px dashed rgba(255, 255, 255, 0.1) !important;
  display: flex;
  align-items: center;
  justify-content: center;
  aspect-ratio: 2/3;
  transition: all 0.25s;
}

.bmol-more-card:hover {
  background: rgba(225, 37, 27, 0.02) !important;
  border-color: rgba(225, 37, 27, 0.35) !important;
}

.more-card-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 100%;
  width: 100%;
}

.more-count {
  font-size: 1.1rem;
  font-weight: 750;
  color: #e1251b;
}

.more-sub {
  font-size: 0.68rem;
  color: #71717a;
}

/* View Header area */
.bmol-view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding-bottom: 1rem;
}

.bmol-section-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
}

.bmol-view-selector {
  display: flex;
  background: rgba(20, 20, 22, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 10px;
  padding: 2px;
}

.view-selector-btn {
  background: transparent;
  border: none;
  color: #9ca3af;
  padding: 0.375rem 0.875rem;
  font-size: 0.8125rem;
  font-weight: 600;
  cursor: pointer;
  border-radius: 8px;
  transition: background-color 0.2s, color 0.2s;
}

.view-selector-btn--active {
  background: rgba(225, 37, 27, 0.15);
  color: #e1251b;
}

/* Rank groups */
.bmol-ranks-list {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.bmol-rank-group {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.bmol-rank-header {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.rank-number-tag {
  font-size: 0.875rem;
  font-weight: 700;
  color: #e1251b;
  background: rgba(225, 37, 27, 0.1);
  border: 1px solid rgba(225, 37, 27, 0.2);
  padding: 4px 10px;
  border-radius: 8px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.rank-line {
  flex: 1;
  height: 1px;
  background: rgba(255, 255, 255, 0.05);
}

/* Scrollable row of items in same rank */
.bmol-rank-items-row {
  display: flex;
  gap: 1.5rem;
  flex-wrap: nowrap;
  overflow: hidden;
  width: 100%;
}

.bmol-rank-items-row--scrollable {
  overflow-x: auto;
  padding-bottom: 0.75rem;
}

.bmol-rank-items-row--scrollable::-webkit-scrollbar {
  height: 6px;
}
.bmol-rank-items-row--scrollable::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 3px;
}
.bmol-rank-items-row--scrollable::-webkit-scrollbar-thumb {
  background: rgba(225, 37, 27, 0.3);
  border-radius: 3px;
}
.bmol-rank-items-row--scrollable::-webkit-scrollbar-thumb:hover {
  background: rgba(225, 37, 27, 0.5);
}

.bmol-rank-items-row .bmol-item-card {
  width: calc(20% - 1.2rem);
  min-width: 150px;
  flex-shrink: 0;
}

.btn-toggle-expand {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  color: #9ca3af;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 4px 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-toggle-expand:hover {
  background: rgba(255, 255, 255, 0.04);
  color: #ffffff;
  border-color: rgba(255, 255, 255, 0.15);
}


.bmol-item-card {
  background: rgba(20, 20, 22, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 16px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: transform 0.25s, border-color 0.25s;
}

.bmol-item-card:hover {
  transform: translateY(-4px);
  border-color: rgba(255, 255, 255, 0.1);
}

.bmol-poster-frame {
  aspect-ratio: 2/3;
  background: #1c1c1e;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.bmol-poster-frame img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.bmol-item-card:hover .bmol-poster-frame img {
  transform: scale(1.05);
}

.bmol-meta {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  flex: 1;
}

.bmol-title {
  font-size: 0.875rem;
  font-weight: 700;
  color: #ffffff;
  text-decoration: none;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s;
  flex: 1;
}

.bmol-title:hover {
  color: #e1251b;
}

.bmol-rating {
  font-size: 0.75rem;
  color: #ffb800;
  font-weight: 600;
}

.bmol-actions {
  display: flex;
  gap: 0.375rem;
  margin-top: 0.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.04);
  padding-top: 0.5rem;
}

.bmol-action-btn {
  flex: 1;
  background: #27272a;
  border: 1px solid rgba(255, 255, 255, 0.05);
  color: #a1a1aa;
  font-size: 0.75rem;
  padding: 4px 0;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s, color 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.bmol-action-btn:hover:not(:disabled) {
  background: #3f3f46;
  color: #ffffff;
}

.bmol-action-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.bmol-action-btn--danger:hover {
  background: rgba(239, 68, 68, 0.1) !important;
  color: #ef4444 !important;
  border-color: rgba(239, 68, 68, 0.2) !important;
}

.small-ring {
  width: 24px;
  height: 24px;
}

/* Filters styling */
.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  background: rgba(20, 20, 22, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.03);
  border-radius: 16px;
  padding: 1rem;
  margin-bottom: 2rem;
  align-items: center;
}

.filter-input-wrapper {
  display: flex;
  align-items: center;
  background: #161618;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  padding: 0 1rem;
  height: 42px;
}

.search-wrapper {
  flex: 1;
  min-width: 200px;
}

.filter-icon, .select-icon {
  color: #9ca3af;
}

.filter-input {
  background: transparent;
  border: none;
  color: #ffffff;
  padding-left: 0.75rem;
  font-size: 0.875rem;
  width: 100%;
  outline: none;
}

.filter-select-wrapper {
  display: flex;
  align-items: center;
  background: #161618;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  padding: 0 1rem;
  height: 42px;
  min-width: 150px;
  position: relative;
}

.filter-select {
  background: transparent;
  border: none;
  color: #ffffff;
  padding-left: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  width: 100%;
  outline: none;
  cursor: pointer;
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
}

.filter-select option {
  background: #161618;
  color: #ffffff;
}

/* Grid lists and cards */
.media-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

@media (min-width: 640px) {
  .media-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (min-width: 768px) {
  .media-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (min-width: 1024px) {
  .media-grid {
    grid-template-columns: repeat(5, 1fr);
  }
}

.media-library-card {
  background: rgba(20, 20, 22, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 16px;
  overflow: hidden;
  transition: transform 0.25s, border-color 0.25s;
  display: flex;
  flex-direction: column;
}

.media-library-card:hover {
  transform: translateY(-4px);
  border-color: rgba(255, 255, 255, 0.12);
}

.media-poster-frame {
  position: relative;
  aspect-ratio: 2/3;
  background: #1c1c1e;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.media-poster-frame img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.media-library-card:hover .media-poster-frame img {
  transform: scale(1.05);
}

.media-type-chip {
  position: absolute;
  top: 8px;
  right: 8px;
  background: rgba(0, 0, 0, 0.75);
  color: #ffffff;
  font-size: 0.625rem;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.media-card-meta {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  flex: 1;
}

.media-title {
  font-size: 0.875rem;
  font-weight: 700;
  color: #ffffff;
  text-decoration: none;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s;
}

.media-title:hover {
  color: #e1251b;
}

.media-details-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.75rem;
}

.media-rating-badge {
  color: #ffb800;
  font-weight: 600;
}

.save-date {
  color: #71717a;
}

.tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.tag-badge {
  font-size: 0.675rem;
  background: rgba(255, 255, 255, 0.05);
  color: #a1a1aa;
  padding: 2px 6px;
  border-radius: 4px;
}

/* Reviews List card styling */
.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.review-library-card {
  background: rgba(20, 20, 22, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 20px;
  padding: 1.5rem;
  display: flex;
  gap: 1.5rem;
}

.review-media-poster {
  width: 80px;
  height: 120px;
  border-radius: 12px;
  overflow: hidden;
  flex-shrink: 0;
  background: #1c1c1e;
  display: flex;
  align-items: center;
  justify-content: center;
}

.review-media-poster img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.review-card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.review-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.review-media-title {
  font-size: 1.125rem;
  font-weight: 700;
  color: #ffffff;
  text-decoration: none;
  transition: color 0.2s;
}

.review-media-title:hover {
  color: #e1251b;
}

.review-card-stars {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  background: rgba(255, 184, 0, 0.08);
  border: 1px solid rgba(255, 184, 0, 0.15);
  padding: 4px 10px;
  border-radius: 8px;
}

.rating-num {
  font-size: 0.8125rem;
  font-weight: 700;
  color: #ffb800;
}

.stars-icon {
  color: #ffb800;
  font-size: 0.8125rem;
}

.review-card-body {
  font-size: 0.875rem;
  color: #d1d5db;
  line-height: 1.6;
  margin: 0;
  white-space: pre-wrap;
}

.review-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
  padding-top: 0.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.04);
}

.review-date {
  font-size: 0.75rem;
  color: #71717a;
}

.review-meta-actions {
  display: flex;
  gap: 1rem;
}

.action-count {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.75rem;
  color: #a1a1aa;
}

.liked-heart {
  color: #ef4444;
}

/* Empty state styling */
.empty-state-card {
  grid-column: 1 / -1;
  text-align: center;
  padding: 4rem 2rem;
  background: rgba(20, 20, 22, 0.3);
  border: 1px dashed rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  color: #71717a;
}

.empty-icon {
  margin-bottom: 1rem;
  color: #3f3f46;
}

/* Poster fallback */
.poster-fallback {
  color: #3f3f46;
}

/* entrance animation */
.fade-in-up {
  opacity: 0;
  transform: translateY(15px);
  animation: fadeInUp 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes fadeInUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Premium Breadcrumb Header Card */
.bmol-breadcrumb-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 14px;
  padding: 12px 18px;
  margin-bottom: 20px;
  backdrop-filter: blur(10px);
}

.breadcrumb-left {
  display: flex;
  align-items: center;
  gap: 14px;
}

.btn-back-circle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #a1a1aa;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.btn-back-circle:hover {
  background: rgba(225, 37, 27, 0.1);
  border-color: rgba(225, 37, 27, 0.3);
  color: #e1251b;
  transform: translateX(-2px);
}

.breadcrumb-path {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
  font-weight: 500;
}

.path-link {
  color: #8a8a8e;
  cursor: pointer;
  transition: color 0.2s;
}

.path-link:hover {
  color: #ffffff;
}

.path-separator {
  color: rgba(255, 255, 255, 0.15);
  font-weight: 300;
}

.path-current {
  display: flex;
  align-items: center;
}

.rank-badge-glow {
  background: linear-gradient(135deg, #e1251b 0%, #b81d15 100%);
  color: #ffffff;
  font-size: 0.78rem;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 6px;
  box-shadow: 0 0 12px rgba(225, 37, 27, 0.35);
  letter-spacing: 0.5px;
}

.detail-count-chip {
  font-size: 0.75rem;
  font-weight: 600;
  color: #8a8a8e;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.05);
  padding: 4px 10px;
  border-radius: 20px;
}

/* Detail View Styles */
.bmol-detail-view {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* --- Detail Toolbar --- */
.detail-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  background: rgba(20, 20, 22, 0.55);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 14px;
  backdrop-filter: blur(12px);
}

.detail-search-box {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
  gap: 8px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  padding: 0 12px;
  height: 38px;
  transition: border-color 0.2s, background-color 0.2s;
}

.detail-search-box:focus-within {
  border-color: rgba(225, 37, 27, 0.35);
  background: rgba(255, 255, 255, 0.04);
}

.detail-search-icon {
  color: #52525b;
  flex-shrink: 0;
  transition: color 0.2s;
}

.detail-search-box:focus-within .detail-search-icon {
  color: #e1251b;
}

.detail-search-input {
  background: transparent;
  border: none;
  outline: none;
  color: #e4e4e7;
  font-size: 0.82rem;
  font-weight: 500;
  width: 100%;
  min-width: 0;
}

.detail-search-input::placeholder {
  color: #52525b;
  font-weight: 400;
}

.detail-search-clear {
  background: rgba(255, 255, 255, 0.06);
  border: none;
  color: #71717a;
  font-size: 0.65rem;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.2s;
}

.detail-search-clear:hover {
  background: rgba(239, 68, 68, 0.12);
  color: #ef4444;
}

.detail-toolbar-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.detail-result-count {
  font-size: 0.72rem;
  font-weight: 600;
  color: #52525b;
  white-space: nowrap;
  letter-spacing: 0.01em;
}

@media (max-width: 575px) {
  .detail-toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
    padding: 10px 12px;
  }
  .detail-toolbar-actions {
    justify-content: space-between;
  }
}

.bmol-detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.25rem;
}
@media (min-width: 576px) {
  .bmol-detail-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
@media (min-width: 768px) {
  .bmol-detail-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
.bmol-detail-grid .bmol-item-card {
  width: 100% !important;
  min-width: 0 !important;
}

/* Modal Fade Transition */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.22s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

/* SPOTLIGHT QUICK ADD STYLING */
.spotlight-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(12px);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  z-index: 99999;
  padding-top: 10vh;
}

.spotlight-modal {
  width: 90%;
  max-width: 600px;
  background: #121214;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  box-shadow: 0 32px 64px rgba(0, 0, 0, 0.8);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: 75vh;
}

/* Selecting Zone (Top) */
.spotlight-selecting-zone {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 14px 18px;
  background: rgba(255, 255, 255, 0.01);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  max-height: 120px;
  overflow-y: auto;
}

.selecting-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(225, 37, 27, 0.08);
  border: 1px solid rgba(225, 37, 27, 0.2);
  border-radius: 20px;
  padding: 4px 10px;
}

.selecting-chip-poster {
  width: 16px;
  height: 24px;
  border-radius: 2px;
  object-fit: cover;
}

.selecting-chip-title {
  font-size: 0.78rem;
  font-weight: 600;
  color: #ffffff;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.btn-remove-selecting {
  background: transparent;
  border: none;
  color: #8a8a8e;
  font-size: 0.72rem;
  cursor: pointer;
  padding: 0 2px;
  transition: color 0.2s;
}

.btn-remove-selecting:hover {
  color: #ef4444;
}

/* Spotlight Search Bar */
.spotlight-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 18px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.spotlight-search-wrapper {
  display: flex;
  align-items: center;
  flex: 1;
  background: #1c1c1e;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  padding: 0 12px;
  height: 42px;
}

.spotlight-search-icon {
  color: #71717a;
}

.spotlight-search-input {
  background: transparent;
  border: none;
  color: #ffffff;
  font-size: 0.9rem;
  width: 100%;
  outline: none;
  padding-left: 10px;
}

.spotlight-toggle-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
}

.switch-label {
  font-size: 0.78rem;
  font-weight: 600;
  color: #8a8a8e;
}

/* Switch styling */
.switch-container {
  position: relative;
  display: inline-block;
  width: 34px;
  height: 20px;
}

.switch-container input {
  opacity: 0;
  width: 0;
  height: 0;
}

.switch-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #27272a;
  transition: .25s;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.switch-slider:before {
  position: absolute;
  content: "";
  height: 12px;
  width: 12px;
  left: 3px;
  bottom: 3px;
  background-color: #8a8a8e;
  transition: .25s;
  border-radius: 50%;
}

.switch-container input:checked + .switch-slider {
  background-color: rgba(225, 37, 27, 0.15);
  border-color: rgba(225, 37, 27, 0.3);
}

.switch-container input:checked + .switch-slider:before {
  transform: translateX(14px);
  background-color: #e1251b;
}

/* Search Results */
.spotlight-results {
  flex: 1;
  overflow-y: auto;
  min-height: 150px;
  max-height: 320px;
}

.spotlight-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30px 0;
}

.spotlight-results-list {
  display: flex;
  flex-direction: column;
}

.spotlight-result-row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 18px;
  cursor: pointer;
  transition: background-color 0.2s;
  border-bottom: 1px solid rgba(255, 255, 255, 0.02);
}

.spotlight-result-row:hover {
  background: rgba(255, 255, 255, 0.03);
}

.spotlight-result-row--selected {
  background: rgba(225, 37, 27, 0.02) !important;
}

.result-row-poster {
  width: 32px;
  height: 48px;
  border-radius: 4px;
  object-fit: cover;
}

.result-row-poster-fallback {
  width: 32px;
  height: 48px;
  border-radius: 4px;
  background: #1c1c1e;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #71717a;
}

.result-row-info {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.result-row-title {
  font-size: 0.85rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
}

.result-row-date {
  font-size: 0.72rem;
  color: #71717a;
  margin: 4px 0 0 0;
}

.result-row-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.2s;
}

.spotlight-result-row:hover .result-row-indicator {
  border-color: rgba(225, 37, 27, 0.5);
  background: rgba(225, 37, 27, 0.1);
}

.indicator-check {
  color: #ef4444;
  font-weight: 700;
  font-size: 0.8rem;
}

.indicator-plus {
  color: #8a8a8e;
  font-size: 0.9rem;
}

.spotlight-result-row:hover .indicator-plus {
  color: #e1251b;
}

.spotlight-empty-results {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  color: #71717a;
  font-size: 0.85rem;
}

/* Spotlight Footer */
.spotlight-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding: 14px 18px;
  background: rgba(255, 255, 255, 0.01);
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.btn-spotlight-cancel {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #a1a1aa;
  padding: 8px 16px;
  font-size: 0.8rem;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-spotlight-cancel:hover {
  background: rgba(255, 255, 255, 0.04);
  color: #ffffff;
}

.btn-spotlight-save {
  background: #e1251b;
  color: #ffffff;
  border: none;
  padding: 8px 16px;
  font-size: 0.8rem;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-spotlight-save:hover:not(:disabled) {
  background: #b81d15;
}

.btn-spotlight-save:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
