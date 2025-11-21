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

-- サンプルデータ: ユーザー (150件)
INSERT INTO users VALUES (1, "Alice Smith", "alice.smith@example.com", 30, "2024-01-01T10:00:00.000");
INSERT INTO users VALUES (2, "Bob Johnson", "bob.johnson@example.com", 25, "2024-01-02T11:00:00.000");
INSERT INTO users VALUES (3, "Charlie Brown", "charlie.brown@example.com", 35, "2024-01-03T12:00:00.000");
INSERT INTO users VALUES (4, "David Wilson", "david.wilson@example.com", 28, "2024-01-04T13:00:00.000");
INSERT INTO users VALUES (5, "Emma Davis", "emma.davis@example.com", 32, "2024-01-05T14:00:00.000");
INSERT INTO users VALUES (6, "Frank Miller", "frank.miller@example.com", 45, "2024-01-06T15:00:00.000");
INSERT INTO users VALUES (7, "Grace Lee", "grace.lee@example.com", 27, "2024-01-07T16:00:00.000");
INSERT INTO users VALUES (8, "Henry Taylor", "henry.taylor@example.com", 38, "2024-01-08T17:00:00.000");
INSERT INTO users VALUES (9, "Ivy Chen", "ivy.chen@example.com", 29, "2024-01-09T18:00:00.000");
INSERT INTO users VALUES (10, "Jack Thompson", "jack.thompson@example.com", 33, "2024-01-10T19:00:00.000");
INSERT INTO users VALUES (11, "Kate Anderson", "kate.anderson@example.com", 26, "2024-01-11T10:00:00.000");
INSERT INTO users VALUES (12, "Leo Martinez", "leo.martinez@example.com", 41, "2024-01-12T11:00:00.000");
INSERT INTO users VALUES (13, "Mia Garcia", "mia.garcia@example.com", 24, "2024-01-13T12:00:00.000");
INSERT INTO users VALUES (14, "Noah Rodriguez", "noah.rodriguez@example.com", 37, "2024-01-14T13:00:00.000");
INSERT INTO users VALUES (15, "Olivia Lopez", "olivia.lopez@example.com", 31, "2024-01-15T14:00:00.000");
INSERT INTO users VALUES (16, "Paul White", "paul.white@example.com", 42, "2024-01-16T15:00:00.000");
INSERT INTO users VALUES (17, "Quinn Harris", "quinn.harris@example.com", 28, "2024-01-17T16:00:00.000");
INSERT INTO users VALUES (18, "Rose Clark", "rose.clark@example.com", 34, "2024-01-18T17:00:00.000");
INSERT INTO users VALUES (19, "Sam Lewis", "sam.lewis@example.com", 29, "2024-01-19T18:00:00.000");
INSERT INTO users VALUES (20, "Tina Walker", "tina.walker@example.com", 36, "2024-01-20T19:00:00.000");
INSERT INTO users VALUES (21, "Uma Young", "uma.young@example.com", 27, "2024-01-21T10:00:00.000");
INSERT INTO users VALUES (22, "Victor King", "victor.king@example.com", 39, "2024-01-22T11:00:00.000");
INSERT INTO users VALUES (23, "Wendy Scott", "wendy.scott@example.com", 31, "2024-01-23T12:00:00.000");
INSERT INTO users VALUES (24, "Xavier Green", "xavier.green@example.com", 44, "2024-01-24T13:00:00.000");
INSERT INTO users VALUES (25, "Yara Adams", "yara.adams@example.com", 26, "2024-01-25T14:00:00.000");
INSERT INTO users VALUES (26, "Zack Baker", "zack.baker@example.com", 33, "2024-01-26T15:00:00.000");
INSERT INTO users VALUES (27, "Amy Nelson", "amy.nelson@example.com", 29, "2024-01-27T16:00:00.000");
INSERT INTO users VALUES (28, "Brian Carter", "brian.carter@example.com", 40, "2024-01-28T17:00:00.000");
INSERT INTO users VALUES (29, "Chloe Mitchell", "chloe.mitchell@example.com", 25, "2024-01-29T18:00:00.000");
INSERT INTO users VALUES (30, "Daniel Perez", "daniel.perez@example.com", 36, "2024-01-30T19:00:00.000");
INSERT INTO users VALUES (31, "Elena Roberts", "elena.roberts@example.com", 32, "2024-01-31T10:00:00.000");
INSERT INTO users VALUES (32, "Felix Turner", "felix.turner@example.com", 28, "2024-02-01T11:00:00.000");
INSERT INTO users VALUES (33, "Gina Phillips", "gina.phillips@example.com", 34, "2024-02-02T12:00:00.000");
INSERT INTO users VALUES (34, "Hugo Campbell", "hugo.campbell@example.com", 41, "2024-02-03T13:00:00.000");
INSERT INTO users VALUES (35, "Iris Parker", "iris.parker@example.com", 27, "2024-02-04T14:00:00.000");
INSERT INTO users VALUES (36, "James Evans", "james.evans@example.com", 38, "2024-02-05T15:00:00.000");
INSERT INTO users VALUES (37, "Kelly Edwards", "kelly.edwards@example.com", 30, "2024-02-06T16:00:00.000");
INSERT INTO users VALUES (38, "Lucas Collins", "lucas.collins@example.com", 35, "2024-02-07T17:00:00.000");
INSERT INTO users VALUES (39, "Maya Stewart", "maya.stewart@example.com", 26, "2024-02-08T18:00:00.000");
INSERT INTO users VALUES (40, "Nathan Morris", "nathan.morris@example.com", 43, "2024-02-09T19:00:00.000");
INSERT INTO users VALUES (41, "Oscar Rogers", "oscar.rogers@example.com", 31, "2024-02-10T10:00:00.000");
INSERT INTO users VALUES (42, "Piper Reed", "piper.reed@example.com", 29, "2024-02-11T11:00:00.000");
INSERT INTO users VALUES (43, "Quentin Cook", "quentin.cook@example.com", 37, "2024-02-12T12:00:00.000");
INSERT INTO users VALUES (44, "Ruby Morgan", "ruby.morgan@example.com", 28, "2024-02-13T13:00:00.000");
INSERT INTO users VALUES (45, "Simon Bell", "simon.bell@example.com", 42, "2024-02-14T14:00:00.000");
INSERT INTO users VALUES (46, "Tara Murphy", "tara.murphy@example.com", 33, "2024-02-15T15:00:00.000");
INSERT INTO users VALUES (47, "Ulysses Bailey", "ulysses.bailey@example.com", 39, "2024-02-16T16:00:00.000");
INSERT INTO users VALUES (48, "Vera Rivera", "vera.rivera@example.com", 25, "2024-02-17T17:00:00.000");
INSERT INTO users VALUES (49, "Wade Cooper", "wade.cooper@example.com", 36, "2024-02-18T18:00:00.000");
INSERT INTO users VALUES (50, "Xena Richardson", "xena.richardson@example.com", 30, "2024-02-19T19:00:00.000");
INSERT INTO users VALUES (51, "Yale Cox", "yale.cox@example.com", 34, "2024-02-20T10:00:00.000");
INSERT INTO users VALUES (52, "Zelda Howard", "zelda.howard@example.com", 27, "2024-02-21T11:00:00.000");
INSERT INTO users VALUES (53, "Aaron Ward", "aaron.ward@example.com", 40, "2024-02-22T12:00:00.000");
INSERT INTO users VALUES (54, "Beth Torres", "beth.torres@example.com", 32, "2024-02-23T13:00:00.000");
INSERT INTO users VALUES (55, "Carl Peterson", "carl.peterson@example.com", 38, "2024-02-24T14:00:00.000");
INSERT INTO users VALUES (56, "Diana Gray", "diana.gray@example.com", 29, "2024-02-25T15:00:00.000");
INSERT INTO users VALUES (57, "Eric Ramirez", "eric.ramirez@example.com", 41, "2024-02-26T16:00:00.000");
INSERT INTO users VALUES (58, "Fiona James", "fiona.james@example.com", 26, "2024-02-27T17:00:00.000");
INSERT INTO users VALUES (59, "George Watson", "george.watson@example.com", 35, "2024-02-28T18:00:00.000");
INSERT INTO users VALUES (60, "Hannah Brooks", "hannah.brooks@example.com", 31, "2024-02-29T19:00:00.000");
INSERT INTO users VALUES (61, "Ian Kelly", "ian.kelly@example.com", 28, "2024-03-01T10:00:00.000");
INSERT INTO users VALUES (62, "Jane Sanders", "jane.sanders@example.com", 33, "2024-03-02T11:00:00.000");
INSERT INTO users VALUES (63, "Kyle Price", "kyle.price@example.com", 37, "2024-03-03T12:00:00.000");
INSERT INTO users VALUES (64, "Laura Bennett", "laura.bennett@example.com", 30, "2024-03-04T13:00:00.000");
INSERT INTO users VALUES (65, "Mark Wood", "mark.wood@example.com", 42, "2024-03-05T14:00:00.000");
INSERT INTO users VALUES (66, "Nina Barnes", "nina.barnes@example.com", 25, "2024-03-06T15:00:00.000");
INSERT INTO users VALUES (67, "Owen Ross", "owen.ross@example.com", 39, "2024-03-07T16:00:00.000");
INSERT INTO users VALUES (68, "Paula Henderson", "paula.henderson@example.com", 34, "2024-03-08T17:00:00.000");
INSERT INTO users VALUES (69, "Quincy Coleman", "quincy.coleman@example.com", 27, "2024-03-09T18:00:00.000");
INSERT INTO users VALUES (70, "Rachel Jenkins", "rachel.jenkins@example.com", 36, "2024-03-10T19:00:00.000");
INSERT INTO users VALUES (71, "Steve Perry", "steve.perry@example.com", 40, "2024-03-11T10:00:00.000");
INSERT INTO users VALUES (72, "Tracy Powell", "tracy.powell@example.com", 29, "2024-03-12T11:00:00.000");
INSERT INTO users VALUES (73, "Ursula Long", "ursula.long@example.com", 32, "2024-03-13T12:00:00.000");
INSERT INTO users VALUES (74, "Vincent Patterson", "vincent.patterson@example.com", 38, "2024-03-14T13:00:00.000");
INSERT INTO users VALUES (75, "Wilma Hughes", "wilma.hughes@example.com", 31, "2024-03-15T14:00:00.000");
INSERT INTO users VALUES (76, "Xander Flores", "xander.flores@example.com", 35, "2024-03-16T15:00:00.000");
INSERT INTO users VALUES (77, "Yasmin Washington", "yasmin.washington@example.com", 26, "2024-03-17T16:00:00.000");
INSERT INTO users VALUES (78, "Zachary Butler", "zachary.butler@example.com", 43, "2024-03-18T17:00:00.000");
INSERT INTO users VALUES (79, "Abigail Simmons", "abigail.simmons@example.com", 28, "2024-03-19T18:00:00.000");
INSERT INTO users VALUES (80, "Bradley Foster", "bradley.foster@example.com", 37, "2024-03-20T19:00:00.000");
INSERT INTO users VALUES (81, "Carmen Gonzales", "carmen.gonzales@example.com", 30, "2024-03-21T10:00:00.000");
INSERT INTO users VALUES (82, "Derek Bryant", "derek.bryant@example.com", 41, "2024-03-22T11:00:00.000");
INSERT INTO users VALUES (83, "Eliza Alexander", "eliza.alexander@example.com", 25, "2024-03-23T12:00:00.000");
INSERT INTO users VALUES (84, "Fernando Russell", "fernando.russell@example.com", 39, "2024-03-24T13:00:00.000");
INSERT INTO users VALUES (85, "Gabriela Griffin", "gabriela.griffin@example.com", 33, "2024-03-25T14:00:00.000");
INSERT INTO users VALUES (86, "Harold Diaz", "harold.diaz@example.com", 36, "2024-03-26T15:00:00.000");
INSERT INTO users VALUES (87, "Isabella Hayes", "isabella.hayes@example.com", 27, "2024-03-27T16:00:00.000");
INSERT INTO users VALUES (88, "Julian Myers", "julian.myers@example.com", 42, "2024-03-28T17:00:00.000");
INSERT INTO users VALUES (89, "Kimberly Ford", "kimberly.ford@example.com", 29, "2024-03-29T18:00:00.000");
INSERT INTO users VALUES (90, "Leonard Hamilton", "leonard.hamilton@example.com", 40, "2024-03-30T19:00:00.000");
INSERT INTO users VALUES (91, "Monica Graham", "monica.graham@example.com", 31, "2024-03-31T10:00:00.000");
INSERT INTO users VALUES (92, "Nicholas Sullivan", "nicholas.sullivan@example.com", 34, "2024-04-01T11:00:00.000");
INSERT INTO users VALUES (93, "Ophelia Wallace", "ophelia.wallace@example.com", 28, "2024-04-02T12:00:00.000");
INSERT INTO users VALUES (94, "Patrick Woods", "patrick.woods@example.com", 38, "2024-04-03T13:00:00.000");
INSERT INTO users VALUES (95, "Queenie Cole", "queenie.cole@example.com", 26, "2024-04-04T14:00:00.000");
INSERT INTO users VALUES (96, "Raymond West", "raymond.west@example.com", 44, "2024-04-05T15:00:00.000");
INSERT INTO users VALUES (97, "Sophia Jordan", "sophia.jordan@example.com", 32, "2024-04-06T16:00:00.000");
INSERT INTO users VALUES (98, "Theodore Owens", "theodore.owens@example.com", 35, "2024-04-07T17:00:00.000");
INSERT INTO users VALUES (99, "Uriel Reynolds", "uriel.reynolds@example.com", 29, "2024-04-08T18:00:00.000");
INSERT INTO users VALUES (100, "Violet Fisher", "violet.fisher@example.com", 37, "2024-04-09T19:00:00.000");
INSERT INTO users VALUES (101, "Walter Ellis", "walter.ellis@example.com", 41, "2024-04-10T10:00:00.000");
INSERT INTO users VALUES (102, "Xiomara Gibson", "xiomara.gibson@example.com", 27, "2024-04-11T11:00:00.000");
INSERT INTO users VALUES (103, "Yvonne Hunt", "yvonne.hunt@example.com", 33, "2024-04-12T12:00:00.000");
INSERT INTO users VALUES (104, "Zane Crawford", "zane.crawford@example.com", 39, "2024-04-13T13:00:00.000");
INSERT INTO users VALUES (105, "Adriana Knight", "adriana.knight@example.com", 30, "2024-04-14T14:00:00.000");
INSERT INTO users VALUES (106, "Benjamin Pierce", "benjamin.pierce@example.com", 36, "2024-04-15T15:00:00.000");
INSERT INTO users VALUES (107, "Cecilia Berry", "cecilia.berry@example.com", 25, "2024-04-16T16:00:00.000");
INSERT INTO users VALUES (108, "Dominic Grant", "dominic.grant@example.com", 42, "2024-04-17T17:00:00.000");
INSERT INTO users VALUES (109, "Esmeralda Wells", "esmeralda.wells@example.com", 28, "2024-04-18T18:00:00.000");
INSERT INTO users VALUES (110, "Franklin Webb", "franklin.webb@example.com", 40, "2024-04-19T19:00:00.000");
INSERT INTO users VALUES (111, "Gemma Simpson", "gemma.simpson@example.com", 31, "2024-04-20T10:00:00.000");
INSERT INTO users VALUES (112, "Harrison Stevens", "harrison.stevens@example.com", 34, "2024-04-21T11:00:00.000");
INSERT INTO users VALUES (113, "Imogen Tucker", "imogen.tucker@example.com", 26, "2024-04-22T12:00:00.000");
INSERT INTO users VALUES (114, "Jasper Porter", "jasper.porter@example.com", 38, "2024-04-23T13:00:00.000");
INSERT INTO users VALUES (115, "Kendra Hunter", "kendra.hunter@example.com", 29, "2024-04-24T14:00:00.000");
INSERT INTO users VALUES (116, "Lorenzo Hicks", "lorenzo.hicks@example.com", 43, "2024-04-25T15:00:00.000");
INSERT INTO users VALUES (117, "Magnolia Crawford", "magnolia.crawford@example.com", 32, "2024-04-26T16:00:00.000");
INSERT INTO users VALUES (118, "Nathaniel Henry", "nathaniel.henry@example.com", 35, "2024-04-27T17:00:00.000");
INSERT INTO users VALUES (119, "Octavia Boyd", "octavia.boyd@example.com", 27, "2024-04-28T18:00:00.000");
INSERT INTO users VALUES (120, "Preston Mason", "preston.mason@example.com", 41, "2024-04-29T19:00:00.000");
INSERT INTO users VALUES (121, "Quintessa Morales", "quintessa.morales@example.com", 30, "2024-04-30T10:00:00.000");
INSERT INTO users VALUES (122, "Reginald Kennedy", "reginald.kennedy@example.com", 36, "2024-05-01T11:00:00.000");
INSERT INTO users VALUES (123, "Serena Warren", "serena.warren@example.com", 28, "2024-05-02T12:00:00.000");
INSERT INTO users VALUES (124, "Tobias Dixon", "tobias.dixon@example.com", 39, "2024-05-03T13:00:00.000");
INSERT INTO users VALUES (125, "Unity Marshall", "unity.marshall@example.com", 25, "2024-05-04T14:00:00.000");
INSERT INTO users VALUES (126, "Valentina Fowler", "valentina.fowler@example.com", 37, "2024-05-05T15:00:00.000");
INSERT INTO users VALUES (127, "Wesley Chambers", "wesley.chambers@example.com", 33, "2024-05-06T16:00:00.000");
INSERT INTO users VALUES (128, "Xyla Rice", "xyla.rice@example.com", 31, "2024-05-07T17:00:00.000");
INSERT INTO users VALUES (129, "Yorick Stone", "yorick.stone@example.com", 40, "2024-05-08T18:00:00.000");
INSERT INTO users VALUES (130, "Zara Hanson", "zara.hanson@example.com", 26, "2024-05-09T19:00:00.000");
INSERT INTO users VALUES (131, "Adrian Ortiz", "adrian.ortiz@example.com", 34, "2024-05-10T10:00:00.000");
INSERT INTO users VALUES (132, "Brianna Newman", "brianna.newman@example.com", 29, "2024-05-11T11:00:00.000");
INSERT INTO users VALUES (133, "Clarence Garrett", "clarence.garrett@example.com", 42, "2024-05-12T12:00:00.000");
INSERT INTO users VALUES (134, "Delilah Welch", "delilah.welch@example.com", 27, "2024-05-13T13:00:00.000");
INSERT INTO users VALUES (135, "Edmund Larson", "edmund.larson@example.com", 38, "2024-05-14T14:00:00.000");
INSERT INTO users VALUES (136, "Francesca Frazier", "francesca.frazier@example.com", 32, "2024-05-15T15:00:00.000");
INSERT INTO users VALUES (137, "Gilbert Burke", "gilbert.burke@example.com", 35, "2024-05-16T16:00:00.000");
INSERT INTO users VALUES (138, "Henrietta Hanson", "henrietta.hanson@example.com", 30, "2024-05-17T17:00:00.000");
INSERT INTO users VALUES (139, "Isaiah Day", "isaiah.day@example.com", 41, "2024-05-18T18:00:00.000");
INSERT INTO users VALUES (140, "Josephine Sharp", "josephine.sharp@example.com", 28, "2024-05-19T19:00:00.000");
INSERT INTO users VALUES (141, "Kingston Boone", "kingston.boone@example.com", 36, "2024-05-20T10:00:00.000");
INSERT INTO users VALUES (142, "Lillian Mann", "lillian.mann@example.com", 25, "2024-05-21T11:00:00.000");
INSERT INTO users VALUES (143, "Montgomery Mack", "montgomery.mack@example.com", 39, "2024-05-22T12:00:00.000");
INSERT INTO users VALUES (144, "Nora Williamson", "nora.williamson@example.com", 33, "2024-05-23T13:00:00.000");
INSERT INTO users VALUES (145, "Orlando Stevenson", "orlando.stevenson@example.com", 37, "2024-05-24T14:00:00.000");
INSERT INTO users VALUES (146, "Penelope Ryan", "penelope.ryan@example.com", 31, "2024-05-25T15:00:00.000");
INSERT INTO users VALUES (147, "Quade Spencer", "quade.spencer@example.com", 40, "2024-05-26T16:00:00.000");
INSERT INTO users VALUES (148, "Rosalind Fernandez", "rosalind.fernandez@example.com", 26, "2024-05-27T17:00:00.000");
INSERT INTO users VALUES (149, "Sebastian Dunn", "sebastian.dunn@example.com", 34, "2024-05-28T18:00:00.000");
INSERT INTO users VALUES (150, "Tabitha Castillo", "tabitha.castillo@example.com", 29, "2024-05-29T19:00:00.000");

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
