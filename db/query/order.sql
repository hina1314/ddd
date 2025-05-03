-- name: SaveOrder :one
INSERT INTO "order" (
    order_sn, user_id, hotel_id, merchant_id, total_price, total_number,
    total_pay_ticket, status, created_at, expire_time
)
VALUES (
           @order_sn::VARCHAR, @user_id::BIGINT, @hotel_id::BIGINT, @merchant_id::BIGINT,
           @total_price::DECIMAL, @total_number::INT, @total_pay_ticket::INT, @status::SMALLINT,
           @created_at::TIMESTAMP, @expire_time::TIMESTAMP
       )
RETURNING *;

-- name: SaveOrderRooms :exec
INSERT INTO order_room (
     order_id, room_type_id, room_item_id, price, status, created_at
)
VALUES (
           unnest(@order_id::BIGINT[]),
           unnest(@room_type_id::BIGINT[]),
           unnest(@room_item_id::BIGINT[]),
           unnest(@price::DECIMAL[]),
           unnest(@status::SMALLINT[]),
           unnest(@created_at::TIMESTAMP[])
       );