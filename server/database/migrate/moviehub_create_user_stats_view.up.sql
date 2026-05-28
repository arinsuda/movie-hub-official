-- View: user_stats
-- คำนวณ stats ของ user แบบ real-time จาก tables จริง
-- ไม่พึ่ง denormalized counter บน users table อีกต่อไป
--
-- review_count  : นับ reviews ที่ยังไม่ถูกลบ (soft-delete safe)
-- follower_count: นับ user_follows ที่ user นี้เป็น followee
-- following_count: นับ user_follows ที่ user นี้เป็น follower
CREATE
OR REPLACE VIEW user_stats AS
SELECT
    u.id AS user_id,
    -- จำนวน reviews ที่ user เขียน (ไม่นับที่ถูก soft-delete)
    COUNT(DISTINCT r.id) AS review_count,
    -- จำนวน followers (คนอื่น follow user นี้)
    COUNT(DISTINCT f_in.id) AS follower_count,
    -- จำนวน following (user นี้ follow คนอื่น)
    COUNT(DISTINCT f_out.id) AS following_count
FROM
    users u
    LEFT JOIN reviews r ON r.user_id = u.id
    AND r.deleted_at IS NULL -- ถ้ายังไม่มี user_follows table ให้ comment 2 join นี้ออกก่อน
    -- แล้วค่อย uncomment เมื่อสร้าง follow system แล้ว
    LEFT JOIN user_follows f_in ON f_in.followee_id = u.id
    LEFT JOIN user_follows f_out ON f_out.follower_id = u.id
GROUP BY
    u.id;