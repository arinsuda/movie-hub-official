-- View: user_stats
-- คำนวณ stats ของ user แบบ real-time จาก tables จริง
--
-- review_count   : reviews ที่ user เขียน (soft-delete safe)
-- like_count     : media ที่ user กด like
-- watchlist_count: library_items ที่ list_type = 'watchlist'
-- watched_count  : library_items ที่ list_type = 'watched'
-- follower_count : คนอื่น follow user นี้ (status = accepted)
-- following_count: user นี้ follow คนอื่น (status = accepted)
CREATE OR REPLACE VIEW user_stats AS
SELECT
    u.id AS user_id,

    COUNT(DISTINCT r.id)   AS review_count,
    COUNT(DISTINCT ml.id)  AS like_count,

    COUNT(DISTINCT CASE WHEN li_w.list_type  = 'watchlist' THEN li_w.id END) AS watchlist_count,
    COUNT(DISTINCT CASE WHEN li_wd.list_type = 'watched'   THEN li_wd.id END) AS watched_count,

    COUNT(DISTINCT f_in.id)  AS follower_count,
    COUNT(DISTINCT f_out.id) AS following_count

FROM users u

-- reviews ที่ยังไม่ถูก soft-delete
LEFT JOIN reviews r
    ON r.user_id = u.id
    AND r.deleted_at IS NULL

-- media ที่ user กด like
LEFT JOIN media_likes ml
    ON ml.user_id = u.id
    AND ml.deleted_at IS NULL

-- watchlist
LEFT JOIN library_items li_w
    ON li_w.user_id = u.id
    AND li_w.list_type = 'watchlist'
    AND li_w.deleted_at IS NULL

-- watched
LEFT JOIN library_items li_wd
    ON li_wd.user_id = u.id
    AND li_wd.list_type = 'watched'
    AND li_wd.deleted_at IS NULL

-- followers
LEFT JOIN user_follows f_in
    ON f_in.followee_id = u.id
    AND f_in.status = 'accepted'

-- following
LEFT JOIN user_follows f_out
    ON f_out.follower_id = u.id
    AND f_out.status = 'accepted'

GROUP BY u.id;
