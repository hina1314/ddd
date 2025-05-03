-- name: CheckBookingConflicts :one
SELECT COUNT(*) AS conflict_count
FROM user_plan
WHERE phone = @phone
  AND status < 3 -- 0 待付款 1 待入住 2 入住中
  AND (
      -- Overlapping date ranges: new plan [start, end) conflicts with existing [start_date, end_date)
      (start_date < @end_date AND end_date > @start_date)
  );

-- name: AddUserPlan :exec
INSERT INTO user_plan (
    order_id, room_item_id, phone, name, start_date, end_date, status,
    created_at, updated_at
)
VALUES (
           unnest(@order_id::BIGINT[]),
           unnest(@room_item_id::BIGINT[]),
           unnest(@phone::VARCHAR[]),
           unnest(@name::VARCHAR[]),
           unnest(@start_date::TIMESTAMP[]),
           unnest(@end_date::TIMESTAMP[]),
           unnest(@status::SMALLINT[]),
           unnest(@created_at::TIMESTAMP[]),
           unnest(@updated_at::TIMESTAMP[])
       );