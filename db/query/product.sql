-- name: ListProducts :many
SELECT
    p.id,
    p.name,
    p.description,
    MIN(s.price) AS price,
    p.images,
    p.created_at
FROM product AS p
         JOIN product_sku s ON p.id = s.product_id
GROUP BY p.id
ORDER BY p.created_at DESC
    LIMIT $2 OFFSET $1;



-- name: GetProductWithSkus :one
SELECT
    p.id AS product_id,
    p.name AS product_name,
    p.description,
    p.price AS product_price,
    p.images AS product_images,
    p.created_at AS product_created_at,
    json_agg(
            json_build_object(
                    'id', s.id,
                    'name', s.name,
                    'specs', s.specs,
                    'price', s.price,
                    'images', s.images,
                    'created_at', s.created_at,
                    'stock', st.stock
            )
    ) AS skus
FROM product AS p
         LEFT JOIN product_sku AS s ON p.id = s.product_id
         LEFT JOIN product_sku_stock AS st ON s.id = st.sku_id
WHERE p.id = $1
GROUP BY p.id;



-- name: UpdateSkuStock :exec
UPDATE product_sku_stock
SET stock = $2
WHERE sku_id = $1;


-- name: DecreaseSkuStock :exec
UPDATE product_sku_stock
SET stock = stock - $2
WHERE sku_id = $1 AND stock >= $2;

