
BEGIN;

-- Step 1: Create a temporary view for recent high-value orders
CREATE VIEW recent_high_value_orders AS
SELECT 
    o.order_id,
    o.customer_id,
    o.order_date,
    SUM(oi.quantity * p.price) AS total_amount
FROM orders o
INNER JOIN order_items oi ON o.order_id = oi.order_id
INNER JOIN products p ON oi.product_id = p.product_id
WHERE o.order_date BETWEEN CURRENT_DATE - INTERVAL '30 days' AND CURRENT_DATE
GROUP BY o.order_id, o.customer_id, o.order_date
HAVING SUM(oi.quantity * p.price) > 500;

-- Step 2: Main query to fetch customer and order details
SELECT DISTINCT
    c.customer_id,
    c.first_name,
    c.last_name,
    c.email,
    r.total_amount,
    CASE 
        WHEN r.total_amount > 1000 THEN 'VIP'
        WHEN r.total_amount > 500 THEN 'Premium'
        ELSE 'Standard'
    END AS customer_tier,
    pc.category_name
FROM customers c
JOIN recent_high_value_orders r ON c.customer_id = r.customer_id
LEFT JOIN orders o ON r.order_id = o.order_id
LEFT JOIN order_items oi ON o.order_id = oi.order_id
LEFT JOIN products p ON oi.product_id = p.product_id
LEFT JOIN product_categories pc ON p.category_id = pc.category_id
WHERE c.is_active = TRUE
  AND c.email IS NOT NULL
ORDER BY r.total_amount DESC
LIMIT 100 OFFSET 0;

-- Step 3: Log the report generation
INSERT INTO audit_log (action, performed_by, performed_at, details)
VALUES ('Generate High Value Report', 'admin_user', CURRENT_TIMESTAMP, 'Top 100 high-value customers report generated');

COMMIT;

