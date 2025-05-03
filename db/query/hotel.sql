-- name: FindSkuByID :one
SELECT id, hotel_id, room_type_id, sales_price, refund_status, merchant_id, created_at, updated_at
FROM hotel_sku
WHERE id = $1;

-- name: FindRoomItemIDsByRoomTypeID :many
SELECT id
FROM hotel_room_item
WHERE hotel_id = $1
  AND room_type_id = $2
  AND status = 1;


-- name: FindAvailableRoomItems :many
SELECT DISTINCT ri.id
FROM hotel_room_item ri
WHERE ri.id = ANY(@room_item_ids::BIGINT[])
  AND NOT EXISTS (
    SELECT 1
    FROM hotel_room_date hrd
    WHERE hrd.room_item_id = ri.id
      AND hrd.date >= @start_date::DATE
        AND hrd.date < @end_date::DATE
) LIMIT $1;


-- name: GetPrice :many
SELECT id, hotel_id, sku_id, room_type_id, date,
    market_price, sale_price, ticket_price, ticket_status, created_at, updated_at
FROM hotel_sku_day_price
WHERE
    sku_id = $1
    AND date >= @start_date::DATE AND date < @end_date::DATE
ORDER BY date;

-- name: AddRoomDate :exec
INSERT INTO hotel_room_date (order_id, hotel_id, room_item_id, date, status, created_at, updated_at)
VALUES (
           unnest(@order_id::BIGINT[]),
           unnest(@hotel_id::BIGINT[]),
           unnest(@room_item_id::BIGINT[]),
           unnest(@date::DATE[]),
           unnest(@status::SMALLINT[]),
           unnest(@created_at::TIMESTAMP[]),
           unnest(@updated_at::TIMESTAMP[])
       );
