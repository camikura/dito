-- Oracle NoSQL Database テスト用テーブル

-- 親テーブル: ユーザー
CREATE TABLE IF NOT EXISTS users (
    id INTEGER,
    name STRING,
    email STRING,
    age INTEGER,
    created_at TIMESTAMP(3),
    PRIMARY KEY(id)
);

-- 子テーブル: ユーザーの住所
CREATE TABLE IF NOT EXISTS users.addresses (
    address_id INTEGER,
    type STRING,
    street STRING,
    city STRING,
    state STRING,
    postal_code STRING,
    country STRING,
    is_primary BOOLEAN,
    PRIMARY KEY(address_id)
);

-- 子テーブル: ユーザーの電話番号
CREATE TABLE IF NOT EXISTS users.phones (
    phone_id INTEGER,
    type STRING,
    number STRING,
    is_primary BOOLEAN,
    PRIMARY KEY(phone_id)
);

-- 親テーブル: 商品
CREATE TABLE IF NOT EXISTS products (
    product_id STRING,
    name STRING,
    price DOUBLE,
    category STRING,
    stock INTEGER,
    description STRING,
    PRIMARY KEY(product_id)
);

-- 子テーブル: 商品レビュー
CREATE TABLE IF NOT EXISTS products.reviews (
    review_id INTEGER,
    user_id INTEGER,
    rating INTEGER,
    comment STRING,
    reviewed_at TIMESTAMP(3),
    PRIMARY KEY(review_id)
);

-- 親テーブル: 注文
CREATE TABLE IF NOT EXISTS orders (
    order_id STRING,
    user_id INTEGER,
    total_price DOUBLE,
    status STRING,
    ordered_at TIMESTAMP(3),
    PRIMARY KEY(order_id)
);

-- 子テーブル: 注文明細
CREATE TABLE IF NOT EXISTS orders.items (
    item_id INTEGER,
    product_id STRING,
    product_name STRING,
    quantity INTEGER,
    unit_price DOUBLE,
    subtotal DOUBLE,
    PRIMARY KEY(item_id)
);

-- サンプルデータ: ユーザー (20件)
INSERT INTO users VALUES (1, "Alice", "alice@example.com", 30, "2024-01-01T10:00:00.000");
INSERT INTO users VALUES (2, "Bob", "bob@example.com", 25, "2024-01-02T11:00:00.000");
INSERT INTO users VALUES (3, "Charlie", "charlie@example.com", 35, "2024-01-03T12:00:00.000");
INSERT INTO users VALUES (4, "David", "david@example.com", 28, "2024-01-04T13:00:00.000");
INSERT INTO users VALUES (5, "Emma", "emma@example.com", 32, "2024-01-05T14:00:00.000");
INSERT INTO users VALUES (6, "Frank", "frank@example.com", 45, "2024-01-06T15:00:00.000");
INSERT INTO users VALUES (7, "Grace", "grace@example.com", 27, "2024-01-07T16:00:00.000");
INSERT INTO users VALUES (8, "Henry", "henry@example.com", 38, "2024-01-08T17:00:00.000");
INSERT INTO users VALUES (9, "Ivy", "ivy@example.com", 29, "2024-01-09T18:00:00.000");
INSERT INTO users VALUES (10, "Jack", "jack@example.com", 33, "2024-01-10T19:00:00.000");
INSERT INTO users VALUES (11, "Kate", "kate@example.com", 26, "2024-01-11T10:00:00.000");
INSERT INTO users VALUES (12, "Leo", "leo@example.com", 41, "2024-01-12T11:00:00.000");
INSERT INTO users VALUES (13, "Mia", "mia@example.com", 24, "2024-01-13T12:00:00.000");
INSERT INTO users VALUES (14, "Noah", "noah@example.com", 37, "2024-01-14T13:00:00.000");
INSERT INTO users VALUES (15, "Olivia", "olivia@example.com", 31, "2024-01-15T14:00:00.000");
INSERT INTO users VALUES (16, "Paul", "paul@example.com", 42, "2024-01-16T15:00:00.000");
INSERT INTO users VALUES (17, "Quinn", "quinn@example.com", 28, "2024-01-17T16:00:00.000");
INSERT INTO users VALUES (18, "Rose", "rose@example.com", 34, "2024-01-18T17:00:00.000");
INSERT INTO users VALUES (19, "Sam", "sam@example.com", 29, "2024-01-19T18:00:00.000");
INSERT INTO users VALUES (20, "Tina", "tina@example.com", 36, "2024-01-20T19:00:00.000");

-- サンプルデータ: ユーザーの住所 (30件)
INSERT INTO users.addresses VALUES (1, 1, "home", "123 Main St", "San Francisco", "CA", "94102", "USA", true);
INSERT INTO users.addresses VALUES (1, 2, "work", "456 Market St", "San Francisco", "CA", "94103", "USA", false);
INSERT INTO users.addresses VALUES (2, 1, "home", "789 Oak Ave", "Los Angeles", "CA", "90001", "USA", true);
INSERT INTO users.addresses VALUES (3, 1, "home", "321 Pine St", "Seattle", "WA", "98101", "USA", true);
INSERT INTO users.addresses VALUES (4, 1, "home", "654 Elm St", "Portland", "OR", "97201", "USA", true);
INSERT INTO users.addresses VALUES (5, 1, "home", "987 Maple Ave", "Denver", "CO", "80201", "USA", true);
INSERT INTO users.addresses VALUES (5, 2, "work", "111 Tech Blvd", "Denver", "CO", "80202", "USA", false);
INSERT INTO users.addresses VALUES (6, 1, "home", "222 Lake Dr", "Chicago", "IL", "60601", "USA", true);
INSERT INTO users.addresses VALUES (7, 1, "home", "333 River Rd", "Austin", "TX", "78701", "USA", true);
INSERT INTO users.addresses VALUES (8, 1, "home", "444 Hill St", "Boston", "MA", "02101", "USA", true);
INSERT INTO users.addresses VALUES (9, 1, "home", "555 Bay Ave", "Miami", "FL", "33101", "USA", true);
INSERT INTO users.addresses VALUES (10, 1, "home", "666 Ocean Blvd", "San Diego", "CA", "92101", "USA", true);
INSERT INTO users.addresses VALUES (11, 1, "home", "777 Park Lane", "New York", "NY", "10001", "USA", true);
INSERT INTO users.addresses VALUES (12, 1, "home", "888 Forest Dr", "Atlanta", "GA", "30301", "USA", true);
INSERT INTO users.addresses VALUES (13, 1, "home", "999 Valley Rd", "Phoenix", "AZ", "85001", "USA", true);
INSERT INTO users.addresses VALUES (14, 1, "home", "101 Mountain Ave", "Salt Lake City", "UT", "84101", "USA", true);
INSERT INTO users.addresses VALUES (15, 1, "home", "202 Desert Blvd", "Las Vegas", "NV", "89101", "USA", true);

-- サンプルデータ: ユーザーの電話番号 (25件)
INSERT INTO users.phones VALUES (1, 1, "mobile", "+1-555-0101", true);
INSERT INTO users.phones VALUES (1, 2, "work", "+1-555-0102", false);
INSERT INTO users.phones VALUES (2, 1, "mobile", "+1-555-0201", true);
INSERT INTO users.phones VALUES (3, 1, "mobile", "+1-555-0301", true);
INSERT INTO users.phones VALUES (4, 1, "mobile", "+1-555-0401", true);
INSERT INTO users.phones VALUES (5, 1, "mobile", "+1-555-0501", true);
INSERT INTO users.phones VALUES (5, 2, "work", "+1-555-0502", false);
INSERT INTO users.phones VALUES (6, 1, "mobile", "+1-555-0601", true);
INSERT INTO users.phones VALUES (7, 1, "mobile", "+1-555-0701", true);
INSERT INTO users.phones VALUES (8, 1, "mobile", "+1-555-0801", true);
INSERT INTO users.phones VALUES (9, 1, "mobile", "+1-555-0901", true);
INSERT INTO users.phones VALUES (10, 1, "mobile", "+1-555-1001", true);
INSERT INTO users.phones VALUES (11, 1, "mobile", "+1-555-1101", true);
INSERT INTO users.phones VALUES (12, 1, "mobile", "+1-555-1201", true);
INSERT INTO users.phones VALUES (13, 1, "mobile", "+1-555-1301", true);
INSERT INTO users.phones VALUES (14, 1, "mobile", "+1-555-1401", true);
INSERT INTO users.phones VALUES (15, 1, "mobile", "+1-555-1501", true);

-- サンプルデータ: 商品 (15件)
INSERT INTO products VALUES ("P001", "Laptop", 1200.00, "Electronics", 10, "High-performance laptop");
INSERT INTO products VALUES ("P002", "Mouse", 25.00, "Electronics", 100, "Wireless mouse");
INSERT INTO products VALUES ("P003", "Keyboard", 75.00, "Electronics", 50, "Mechanical keyboard");
INSERT INTO products VALUES ("P004", "Monitor", 300.00, "Electronics", 20, "27-inch 4K monitor");
INSERT INTO products VALUES ("P005", "Headphones", 150.00, "Electronics", 30, "Noise-canceling headphones");
INSERT INTO products VALUES ("P006", "Webcam", 80.00, "Electronics", 40, "HD webcam");
INSERT INTO products VALUES ("P007", "USB Cable", 10.00, "Electronics", 200, "USB-C cable");
INSERT INTO products VALUES ("P008", "Desk Lamp", 45.00, "Furniture", 25, "LED desk lamp");
INSERT INTO products VALUES ("P009", "Office Chair", 250.00, "Furniture", 15, "Ergonomic office chair");
INSERT INTO products VALUES ("P010", "Desk", 400.00, "Furniture", 10, "Standing desk");
INSERT INTO products VALUES ("P011", "Notebook", 5.00, "Stationery", 500, "A4 notebook");
INSERT INTO products VALUES ("P012", "Pen Set", 15.00, "Stationery", 100, "Set of 10 pens");
INSERT INTO products VALUES ("P013", "Backpack", 60.00, "Accessories", 35, "Laptop backpack");
INSERT INTO products VALUES ("P014", "Water Bottle", 20.00, "Accessories", 80, "Insulated water bottle");
INSERT INTO products VALUES ("P015", "Phone Stand", 12.00, "Accessories", 120, "Adjustable phone stand");

-- サンプルデータ: 商品レビュー (25件)
INSERT INTO products.reviews VALUES ("P001", 1, 1, 5, "Excellent laptop! Very fast.", "2024-02-15T10:00:00.000");
INSERT INTO products.reviews VALUES ("P001", 2, 2, 4, "Good performance but a bit heavy.", "2024-02-16T11:00:00.000");
INSERT INTO products.reviews VALUES ("P001", 3, 3, 5, "Best laptop I've ever owned!", "2024-02-17T12:00:00.000");
INSERT INTO products.reviews VALUES ("P002", 1, 1, 5, "Perfect mouse for work.", "2024-02-17T12:00:00.000");
INSERT INTO products.reviews VALUES ("P002", 2, 4, 4, "Good mouse, battery lasts long.", "2024-02-18T13:00:00.000");
INSERT INTO products.reviews VALUES ("P003", 1, 5, 5, "Love the mechanical feel!", "2024-02-19T14:00:00.000");
INSERT INTO products.reviews VALUES ("P003", 2, 6, 3, "Good but a bit noisy.", "2024-02-20T15:00:00.000");
INSERT INTO products.reviews VALUES ("P004", 1, 7, 5, "Amazing picture quality!", "2024-02-21T16:00:00.000");
INSERT INTO products.reviews VALUES ("P005", 1, 8, 5, "Best noise cancellation.", "2024-02-22T17:00:00.000");
INSERT INTO products.reviews VALUES ("P005", 2, 9, 4, "Comfortable for long use.", "2024-02-23T18:00:00.000");
INSERT INTO products.reviews VALUES ("P006", 1, 10, 4, "Clear video quality.", "2024-02-24T10:00:00.000");
INSERT INTO products.reviews VALUES ("P007", 1, 11, 5, "Durable cable.", "2024-02-25T11:00:00.000");
INSERT INTO products.reviews VALUES ("P008", 1, 12, 5, "Perfect lighting for desk.", "2024-02-26T12:00:00.000");
INSERT INTO products.reviews VALUES ("P009", 1, 13, 5, "Very comfortable chair.", "2024-02-27T13:00:00.000");
INSERT INTO products.reviews VALUES ("P009", 2, 14, 4, "Good support for back.", "2024-02-28T14:00:00.000");
INSERT INTO products.reviews VALUES ("P010", 1, 15, 5, "Love the standing feature!", "2024-02-29T15:00:00.000");

-- サンプルデータ: 注文 (15件)
INSERT INTO orders VALUES ("O001", 1, 1275.00, "completed", "2024-02-01T14:00:00.000");
INSERT INTO orders VALUES ("O002", 2, 50.00, "pending", "2024-02-02T15:00:00.000");
INSERT INTO orders VALUES ("O003", 1, 375.00, "shipped", "2024-02-03T16:00:00.000");
INSERT INTO orders VALUES ("O004", 3, 150.00, "completed", "2024-02-04T10:00:00.000");
INSERT INTO orders VALUES ("O005", 4, 80.00, "shipped", "2024-02-05T11:00:00.000");
INSERT INTO orders VALUES ("O006", 5, 460.00, "completed", "2024-02-06T12:00:00.000");
INSERT INTO orders VALUES ("O007", 6, 250.00, "pending", "2024-02-07T13:00:00.000");
INSERT INTO orders VALUES ("O008", 7, 400.00, "shipped", "2024-02-08T14:00:00.000");
INSERT INTO orders VALUES ("O009", 8, 35.00, "completed", "2024-02-09T15:00:00.000");
INSERT INTO orders VALUES ("O010", 9, 1200.00, "pending", "2024-02-10T16:00:00.000");
INSERT INTO orders VALUES ("O011", 10, 92.00, "shipped", "2024-02-11T17:00:00.000");
INSERT INTO orders VALUES ("O012", 11, 325.00, "completed", "2024-02-12T18:00:00.000");
INSERT INTO orders VALUES ("O013", 12, 155.00, "pending", "2024-02-13T19:00:00.000");
INSERT INTO orders VALUES ("O014", 13, 60.00, "shipped", "2024-02-14T10:00:00.000");
INSERT INTO orders VALUES ("O015", 14, 20.00, "completed", "2024-02-15T11:00:00.000");

-- サンプルデータ: 注文明細 (35件)
INSERT INTO orders.items VALUES ("O001", 1, "P001", "Laptop", 1, 1200.00, 1200.00);
INSERT INTO orders.items VALUES ("O001", 2, "P003", "Keyboard", 1, 75.00, 75.00);
INSERT INTO orders.items VALUES ("O002", 1, "P002", "Mouse", 2, 25.00, 50.00);
INSERT INTO orders.items VALUES ("O003", 1, "P004", "Monitor", 1, 300.00, 300.00);
INSERT INTO orders.items VALUES ("O003", 2, "P003", "Keyboard", 1, 75.00, 75.00);
INSERT INTO orders.items VALUES ("O004", 1, "P005", "Headphones", 1, 150.00, 150.00);
INSERT INTO orders.items VALUES ("O005", 1, "P006", "Webcam", 1, 80.00, 80.00);
INSERT INTO orders.items VALUES ("O006", 1, "P009", "Office Chair", 1, 250.00, 250.00);
INSERT INTO orders.items VALUES ("O006", 2, "P007", "USB Cable", 1, 10.00, 10.00);
INSERT INTO orders.items VALUES ("O006", 3, "P001", "Laptop", 1, 1200.00, 1200.00);
INSERT INTO orders.items VALUES ("O007", 1, "P009", "Office Chair", 1, 250.00, 250.00);
INSERT INTO orders.items VALUES ("O008", 1, "P010", "Desk", 1, 400.00, 400.00);
INSERT INTO orders.items VALUES ("O009", 1, "P011", "Notebook", 3, 5.00, 15.00);
INSERT INTO orders.items VALUES ("O009", 2, "P012", "Pen Set", 1, 15.00, 15.00);
INSERT INTO orders.items VALUES ("O009", 3, "P007", "USB Cable", 1, 10.00, 10.00);
INSERT INTO orders.items VALUES ("O010", 1, "P001", "Laptop", 1, 1200.00, 1200.00);
INSERT INTO orders.items VALUES ("O011", 1, "P013", "Backpack", 1, 60.00, 60.00);
INSERT INTO orders.items VALUES ("O011", 2, "P014", "Water Bottle", 1, 20.00, 20.00);
INSERT INTO orders.items VALUES ("O011", 3, "P015", "Phone Stand", 1, 12.00, 12.00);
INSERT INTO orders.items VALUES ("O012", 1, "P004", "Monitor", 1, 300.00, 300.00);
INSERT INTO orders.items VALUES ("O012", 2, "P002", "Mouse", 1, 25.00, 25.00);
INSERT INTO orders.items VALUES ("O013", 1, "P005", "Headphones", 1, 150.00, 150.00);
INSERT INTO orders.items VALUES ("O013", 2, "P007", "USB Cable", 1, 10.00, 10.00);
INSERT INTO orders.items VALUES ("O014", 1, "P013", "Backpack", 1, 60.00, 60.00);
INSERT INTO orders.items VALUES ("O015", 1, "P014", "Water Bottle", 1, 20.00, 20.00);

-- インデックス作成
CREATE INDEX IF NOT EXISTS email_idx ON users (email);
CREATE INDEX IF NOT EXISTS name_idx ON users (name);
CREATE INDEX IF NOT EXISTS category_idx ON products (category);
CREATE INDEX IF NOT EXISTS status_idx ON orders (status);
