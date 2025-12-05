-- Oracle NoSQL Database テスト用テーブル

-- 親テーブル: ユーザー
CREATE TABLE IF NOT EXISTS users (
    id INTEGER,
    name STRING,
    email STRING,
    age INTEGER,
    phone STRING,
    address STRING,
    company STRING,
    job_title STRING,
    department STRING,
    status STRING,
    last_login TIMESTAMP(3),
    notes STRING,
    created_at TIMESTAMP(3),
    preferences JSON,
    metadata JSON,
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
INSERT INTO users VALUES (1, "Alice Smith", "alice.smith@example.com", 30, "+1-415-555-0101", "123 Main St, San Francisco, CA", "Tech Corp", "Software Engineer", "Engineering", "active", "2024-01-01T09:30:00.000", "Experienced developer", "2024-01-01T10:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["developer","senior"],"score":95});
INSERT INTO users VALUES (2, "Bob Johnson", "bob.johnson@example.com", 25, "+1-213-555-0202", "456 Oak Ave, Los Angeles, CA", "Design Studio", "UX Designer", "Product", "active", "2024-01-02T10:15:00.000", "Creative thinker", "2024-01-02T11:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["designer","creative"],"score":88});
INSERT INTO users VALUES (3, "Charlie Brown", "charlie.brown@example.com", 35, "+1-206-555-0303", "789 Pine St, Seattle, WA", "Data Analytics Inc", "Data Scientist", "Analytics", "active", "2024-01-03T11:45:00.000", "ML specialist", "2024-01-03T12:00:00.000", {"theme":"dark","notifications":true,"language":"en","experimental_features":true}, {"tags":["data","machine-learning","ai"],"score":92,"certifications":["AWS","GCP"]});
INSERT INTO users VALUES (4, "David Wilson", "david.wilson@example.com", 28, "+1-503-555-0404", "321 Elm St, Portland, OR", "Cloud Services", "DevOps Engineer", "Infrastructure", "active", "2024-01-04T12:30:00.000", "Cloud expert", "2024-01-04T13:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["devops","cloud"],"score":90});
INSERT INTO users VALUES (5, "Emma Davis", "emma.davis@example.com", 32, "+1-303-555-0505", "654 Maple Ave, Denver, CO", "Marketing Pro", "Marketing Manager", "Marketing", "active", "2024-01-05T13:20:00.000", "Digital marketing guru", "2024-01-05T14:00:00.000", {"theme":"light","notifications":true,"language":"en","email_digest":"weekly"}, {"tags":["marketing","manager"],"score":87,"achievements":["top_performer_2023","innovation_award"]});
INSERT INTO users VALUES (6, "Frank Miller", "frank.miller@example.com", 45, "+1-312-555-0606", "987 Lake Dr, Chicago, IL", "Finance Group", "Financial Analyst", "Finance", "inactive", "2024-01-06T14:10:00.000", "CFA certified", "2024-01-06T15:00:00.000", NULL, NULL);
INSERT INTO users VALUES (7, "Grace Lee", "grace.lee@example.com", 27, "+1-512-555-0707", "111 River Rd, Austin, TX", "Startup Labs", "Product Manager", "Product", "active", "2024-01-07T15:00:00.000", "Agile enthusiast", "2024-01-07T16:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["product","agile"],"score":89});
INSERT INTO users VALUES (8, "Henry Taylor", "henry.taylor@example.com", 38, "+1-617-555-0808", "222 Hill St, Boston, MA", "Consulting Firm", "Senior Consultant", "Strategy", "active", "2024-01-08T16:45:00.000", "MBA from Harvard", "2024-01-08T17:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["consulting","strategy","mba"],"score":93});
INSERT INTO users VALUES (9, "Ivy Chen", "ivy.chen@example.com", 29, "+1-305-555-0909", "333 Bay Ave, Miami, FL", "E-commerce Co", "Frontend Developer", "Engineering", "active", "2024-01-09T17:30:00.000", "React specialist", "2024-01-09T18:00:00.000", {"theme":"dark","notifications":true,"language":"en","beta_features":true}, {"tags":["frontend","react","javascript"],"score":91});
INSERT INTO users VALUES (10, "Jack Thompson", "jack.thompson@example.com", 33, "+1-619-555-1010", "444 Ocean Blvd, San Diego, CA", "Mobile Apps Inc", "Mobile Developer", "Engineering", "active", "2024-01-10T18:15:00.000", "iOS/Android expert", "2024-01-10T19:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["mobile","ios","android"],"score":94});
INSERT INTO users VALUES (11, "Kate Anderson", "kate.anderson@example.com", 26, "+1-212-555-1101", "Address for Kate Anderson", "Company 1", "Engineer", "Engineering", "active", "2024-01-11T09:30:00.000", "User note 11", "2024-01-11T10:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["engineer"],"score":85});
INSERT INTO users VALUES (12, "Leo Martinez", "leo.martinez@example.com", 41, "+1-310-555-1201", "Address for Leo Martinez", "Company 2", "Manager", "Sales", "active", "2024-01-12T10:30:00.000", "User note 12", "2024-01-12T11:00:00.000", {"theme":"dark","notifications":true,"language":"es"}, {"tags":["sales","manager","bilingual"],"score":88});
INSERT INTO users VALUES (13, "Mia Garcia", "mia.garcia@example.com", 24, "+1-415-555-1301", "Address for Mia Garcia", "Company 3", "Designer", "Marketing", "active", "2024-01-13T11:30:00.000", "User note 13", "2024-01-13T12:00:00.000", {"theme":"light","notifications":true,"language":"es"}, {"tags":["design","creative"],"score":82});
INSERT INTO users VALUES (14, "Noah Rodriguez", "noah.rodriguez@example.com", 37, "+1-206-555-1401", "Address for Noah Rodriguez", "Company 4", "Analyst", "Finance", "active", "2024-01-14T12:30:00.000", "User note 14", "2024-01-14T13:00:00.000", {"theme":"auto","notifications":false,"language":"en"}, {"tags":["analyst","finance"],"score":90});
INSERT INTO users VALUES (15, "Olivia Lopez", "olivia.lopez@example.com", 31, "+1-503-555-1501", "Address for Olivia Lopez", "Company 5", "Engineer", "Engineering", "active", "2024-01-15T13:30:00.000", "User note 15", "2024-01-15T14:00:00.000", {"theme":"dark","notifications":true,"language":"en","compact_view":true}, {"tags":["engineer","backend"],"score":89});
INSERT INTO users VALUES (16, "Paul White", "paul.white@example.com", 42, "+1-303-555-1601", "Address for Paul White", "Company 6", "Manager", "Sales", "inactive", "2024-01-16T14:30:00.000", "User note 16", "2024-01-16T15:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["sales"],"score":78});
INSERT INTO users VALUES (17, "Quinn Harris", "quinn.harris@example.com", 28, "+1-312-555-1701", "Address for Quinn Harris", "Company 7", "Designer", "Marketing", "active", "2024-01-17T15:30:00.000", "User note 17", "2024-01-17T16:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["designer","ux"],"score":86});
INSERT INTO users VALUES (18, "Rose Clark", "rose.clark@example.com", 34, "+1-512-555-1801", "Address for Rose Clark", "Company 8", "Analyst", "Finance", "active", "2024-01-18T16:30:00.000", "User note 18", "2024-01-18T17:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["analyst","cpa"],"score":91});
INSERT INTO users VALUES (19, "Sam Lewis", "sam.lewis@example.com", 29, "+1-617-555-1901", "Address for Sam Lewis", "Company 9", "Engineer", "Engineering", "active", "2024-01-19T17:30:00.000", "User note 19", "2024-01-19T18:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["engineer","fullstack"],"score":87});
INSERT INTO users VALUES (20, "Tina Walker", "tina.walker@example.com", 36, "+1-305-555-2001", "Address for Tina Walker", "Company 0", "Manager", "Sales", "active", "2024-01-20T18:30:00.000", "User note 20", "2024-01-20T19:00:00.000", {"theme":"light","notifications":true,"language":"en","email_digest":"daily"}, {"tags":["sales","manager"],"score":84});
INSERT INTO users VALUES (21, "Uma Young", "uma.young@example.com", 27, "+1-619-555-2101", "Address for Uma Young", "Company 1", "Designer", "Marketing", "active", "2024-01-21T09:30:00.000", "User note 21", "2024-01-21T10:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["designer"],"score":83});
INSERT INTO users VALUES (22, "Victor King", "victor.king@example.com", 39, "+1-212-555-2201", "Address for Victor King", "Company 2", "Analyst", "Finance", "active", "2024-01-22T10:30:00.000", "User note 22", "2024-01-22T11:00:00.000", {"theme":"dark","notifications":false,"language":"en","two_factor":true}, {"tags":["analyst","finance","security"],"score":92});
INSERT INTO users VALUES (23, "Wendy Scott", "wendy.scott@example.com", 31, "+1-310-555-2301", "Address for Wendy Scott", "Company 3", "Engineer", "Engineering", "active", "2024-01-23T11:30:00.000", "User note 23", "2024-01-23T12:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["engineer","python"],"score":88});
INSERT INTO users VALUES (24, "Xavier Green", "xavier.green@example.com", 44, "+1-415-555-2401", "Address for Xavier Green", "Company 4", "Manager", "Sales", "active", "2024-01-24T12:30:00.000", "User note 24", "2024-01-24T13:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["sales","manager","leadership"],"score":86});
INSERT INTO users VALUES (25, "Yara Adams", "yara.adams@example.com", 26, "+1-206-555-2501", "Address for Yara Adams", "Company 5", "Designer", "Marketing", "active", "2024-01-25T13:30:00.000", "User note 25", "2024-01-25T14:00:00.000", NULL, NULL);
INSERT INTO users VALUES (26, "Zack Baker", "zack.baker@example.com", 33, "+1-503-555-2601", "Address for Zack Baker", "Company 6", "Analyst", "Finance", "active", "2024-01-26T14:30:00.000", "User note 26", "2024-01-26T15:00:00.000", {"theme":"auto","notifications":false,"language":"en"}, {"tags":["analyst"],"score":84});
INSERT INTO users VALUES (27, "Amy Nelson", "amy.nelson@example.com", 29, "+1-303-555-2701", "Address for Amy Nelson", "Company 7", "Engineer", "Engineering", "active", "2024-01-27T15:30:00.000", "User note 27", "2024-01-27T16:00:00.000", {"theme":"dark","notifications":true,"language":"en","compact_view":true}, {"tags":["engineer","golang"],"score":90});
INSERT INTO users VALUES (28, "Brian Carter", "brian.carter@example.com", 40, "+1-312-555-2801", "Address for Brian Carter", "Company 8", "Manager", "Sales", "active", "2024-01-28T16:30:00.000", "User note 28", "2024-01-28T17:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["sales","manager"],"score":81});
INSERT INTO users VALUES (29, "Chloe Mitchell", "chloe.mitchell@example.com", 25, "+1-512-555-2901", "Address for Chloe Mitchell", "Company 9", "Designer", "Marketing", "active", "2024-01-29T17:30:00.000", "User note 29", "2024-01-29T18:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["designer","ui"],"score":85});
INSERT INTO users VALUES (30, "Daniel Perez", "daniel.perez@example.com", 36, "+1-617-555-3001", "Address for Daniel Perez", "Company 0", "Analyst", "Finance", "active", "2024-01-30T18:30:00.000", "User note 30", "2024-01-30T19:00:00.000", {"theme":"dark","notifications":true,"language":"es"}, {"tags":["analyst","bilingual"],"score":89});
INSERT INTO users VALUES (31, "Elena Roberts", "elena.roberts@example.com", 32, "+1-305-555-3101", "Address for Elena Roberts", "Company 1", "Engineer", "Engineering", "active", "2024-01-31T09:30:00.000", "User note 31", "2024-01-31T10:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["engineer","rust"],"score":92});
INSERT INTO users VALUES (32, "Felix Turner", "felix.turner@example.com", 28, "+1-619-555-3201", "Address for Felix Turner", "Company 2", "Manager", "Sales", "active", "2024-02-01T10:30:00.000", "User note 32", "2024-02-01T11:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["sales"],"score":80});
INSERT INTO users VALUES (33, "Gina Phillips", "gina.phillips@example.com", 34, "+1-212-555-3301", "Address for Gina Phillips", "Company 3", "Designer", "Marketing", "active", "2024-02-02T11:30:00.000", "User note 33", "2024-02-02T12:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["designer","branding"],"score":87});
INSERT INTO users VALUES (34, "Hugo Campbell", "hugo.campbell@example.com", 41, "+1-310-555-3401", "Address for Hugo Campbell", "Company 4", "Analyst", "Finance", "inactive", "2024-02-03T12:30:00.000", "User note 34", "2024-02-03T13:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["analyst"],"score":75});
INSERT INTO users VALUES (35, "Iris Parker", "iris.parker@example.com", 27, "+1-415-555-3501", "Address for Iris Parker", "Company 5", "Engineer", "Engineering", "active", "2024-02-04T13:30:00.000", "User note 35", "2024-02-04T14:00:00.000", {"theme":"dark","notifications":true,"language":"en","beta_features":true}, {"tags":["engineer","typescript"],"score":88});
INSERT INTO users VALUES (36, "James Evans", "james.evans@example.com", 38, "+1-206-555-3601", "Address for James Evans", "Company 6", "Manager", "Sales", "active", "2024-02-05T14:30:00.000", "User note 36", "2024-02-05T15:00:00.000", {"theme":"light","notifications":true,"language":"en","email_digest":"weekly"}, {"tags":["sales","manager"],"score":85});
INSERT INTO users VALUES (37, "Kelly Edwards", "kelly.edwards@example.com", 30, "+1-503-555-3701", "Address for Kelly Edwards", "Company 7", "Designer", "Marketing", "active", "2024-02-06T15:30:00.000", "User note 37", "2024-02-06T16:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["designer"],"score":82});
INSERT INTO users VALUES (38, "Lucas Collins", "lucas.collins@example.com", 35, "+1-303-555-3801", "Address for Lucas Collins", "Company 8", "Analyst", "Finance", "active", "2024-02-07T16:30:00.000", "User note 38", "2024-02-07T17:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["analyst","excel"],"score":86});
INSERT INTO users VALUES (39, "Maya Stewart", "maya.stewart@example.com", 26, "+1-312-555-3901", "Address for Maya Stewart", "Company 9", "Engineer", "Engineering", "active", "2024-02-08T17:30:00.000", "User note 39", "2024-02-08T18:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["engineer","java"],"score":84});
INSERT INTO users VALUES (40, "Nathan Morris", "nathan.morris@example.com", 43, "+1-512-555-4001", "Address for Nathan Morris", "Company 0", "Manager", "Sales", "active", "2024-02-09T18:30:00.000", "User note 40", "2024-02-09T19:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["sales","veteran"],"score":91});
INSERT INTO users VALUES (41, "Oscar Rogers", "oscar.rogers@example.com", 31, "+1-617-555-4101", "Address for Oscar Rogers", "Company 1", "Designer", "Marketing", "active", "2024-02-10T09:30:00.000", "User note 41", "2024-02-10T10:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["designer","figma"],"score":89});
INSERT INTO users VALUES (42, "Piper Reed", "piper.reed@example.com", 29, "+1-305-555-4201", "Address for Piper Reed", "Company 2", "Analyst", "Finance", "active", "2024-02-11T10:30:00.000", "User note 42", "2024-02-11T11:00:00.000", {"theme":"dark","notifications":true,"language":"en","compact_view":true}, {"tags":["analyst","data"],"score":87});
INSERT INTO users VALUES (43, "Quentin Cook", "quentin.cook@example.com", 37, "+1-619-555-4301", "Address for Quentin Cook", "Company 3", "Engineer", "Engineering", "active", "2024-02-12T11:30:00.000", "User note 43", "2024-02-12T12:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["engineer","cpp"],"score":93});
INSERT INTO users VALUES (44, "Ruby Morgan", "ruby.morgan@example.com", 28, "+1-212-555-4401", "Address for Ruby Morgan", "Company 4", "Manager", "Sales", "active", "2024-02-13T12:30:00.000", "User note 44", "2024-02-13T13:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["sales"],"score":83});
INSERT INTO users VALUES (45, "Simon Bell", "simon.bell@example.com", 42, "+1-310-555-4501", "Address for Simon Bell", "Company 5", "Designer", "Marketing", "active", "2024-02-14T13:30:00.000", "User note 45", "2024-02-14T14:00:00.000", NULL, NULL);
INSERT INTO users VALUES (46, "Tara Murphy", "tara.murphy@example.com", 33, "+1-415-555-4601", "Address for Tara Murphy", "Company 6", "Analyst", "Finance", "active", "2024-02-15T14:30:00.000", "User note 46", "2024-02-15T15:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["analyst"],"score":85});
INSERT INTO users VALUES (47, "Ulysses Bailey", "ulysses.bailey@example.com", 39, "+1-206-555-4701", "Address for Ulysses Bailey", "Company 7", "Engineer", "Engineering", "active", "2024-02-16T15:30:00.000", "User note 47", "2024-02-16T16:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["engineer","devops"],"score":88});
INSERT INTO users VALUES (48, "Vera Rivera", "vera.rivera@example.com", 25, "+1-503-555-4801", "Address for Vera Rivera", "Company 8", "Manager", "Sales", "active", "2024-02-17T16:30:00.000", "User note 48", "2024-02-17T17:00:00.000", {"theme":"light","notifications":true,"language":"es"}, {"tags":["sales","bilingual"],"score":81});
INSERT INTO users VALUES (49, "Wade Cooper", "wade.cooper@example.com", 36, "+1-303-555-4901", "Address for Wade Cooper", "Company 9", "Designer", "Marketing", "active", "2024-02-18T17:30:00.000", "User note 49", "2024-02-18T18:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["designer","animation"],"score":90});
INSERT INTO users VALUES (50, "Xena Richardson", "xena.richardson@example.com", 30, "+1-312-555-5001", "Address for Xena Richardson", "Company 0", "Analyst", "Finance", "active", "2024-02-19T18:30:00.000", "User note 50", "2024-02-19T19:00:00.000", {"theme":"auto","notifications":false,"language":"en","two_factor":true}, {"tags":["analyst","security"],"score":92});
INSERT INTO users VALUES (51, "Yale Cox", "yale.cox@example.com", 34, "+1-512-555-5101", "Address for Yale Cox", "Company 1", "Engineer", "Engineering", "active", "2024-02-20T09:30:00.000", "User note 51", "2024-02-20T10:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["engineer"],"score":86});
INSERT INTO users VALUES (52, "Zelda Howard", "zelda.howard@example.com", 27, "+1-617-555-5201", "Address for Zelda Howard", "Company 2", "Manager", "Sales", "active", "2024-02-21T10:30:00.000", "User note 52", "2024-02-21T11:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["sales"],"score":79});
INSERT INTO users VALUES (53, "Aaron Ward", "aaron.ward@example.com", 40, "+1-305-555-5301", "Address for Aaron Ward", "Company 3", "Designer", "Marketing", "active", "2024-02-22T11:30:00.000", "User note 53", "2024-02-22T12:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["designer","senior"],"score":91});
INSERT INTO users VALUES (54, "Beth Torres", "beth.torres@example.com", 32, "+1-619-555-5401", "Address for Beth Torres", "Company 4", "Analyst", "Finance", "active", "2024-02-23T12:30:00.000", "User note 54", "2024-02-23T13:00:00.000", {"theme":"dark","notifications":true,"language":"en","compact_view":true}, {"tags":["analyst"],"score":88});
INSERT INTO users VALUES (55, "Carl Peterson", "carl.peterson@example.com", 38, "+1-212-555-5501", "Address for Carl Peterson", "Company 5", "Engineer", "Engineering", "active", "2024-02-24T13:30:00.000", "User note 55", "2024-02-24T14:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["engineer","kotlin"],"score":87});
INSERT INTO users VALUES (56, "Diana Gray", "diana.gray@example.com", 29, "+1-310-555-5601", "Address for Diana Gray", "Company 6", "Manager", "Sales", "active", "2024-02-25T14:30:00.000", "User note 56", "2024-02-25T15:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["sales"],"score":82});
INSERT INTO users VALUES (57, "Eric Ramirez", "eric.ramirez@example.com", 41, "+1-415-555-5701", "Address for Eric Ramirez", "Company 7", "Designer", "Marketing", "inactive", "2024-02-26T15:30:00.000", "User note 57", "2024-02-26T16:00:00.000", {"theme":"light","notifications":false,"language":"es"}, {"tags":["designer"],"score":76});
INSERT INTO users VALUES (58, "Fiona James", "fiona.james@example.com", 26, "+1-206-555-5801", "Address for Fiona James", "Company 8", "Analyst", "Finance", "active", "2024-02-27T16:30:00.000", "User note 58", "2024-02-27T17:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["analyst"],"score":84});
INSERT INTO users VALUES (59, "George Watson", "george.watson@example.com", 35, "+1-503-555-5901", "Address for George Watson", "Company 9", "Engineer", "Engineering", "active", "2024-02-28T17:30:00.000", "User note 59", "2024-02-28T18:00:00.000", {"theme":"dark","notifications":false,"language":"en","beta_features":true}, {"tags":["engineer","scala"],"score":89});
INSERT INTO users VALUES (60, "Hannah Brooks", "hannah.brooks@example.com", 31, "+1-303-555-6001", "Address for Hannah Brooks", "Company 0", "Manager", "Sales", "active", "2024-02-29T18:30:00.000", "User note 60", "2024-02-29T19:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["sales","manager"],"score":85});
INSERT INTO users VALUES (61, "Ian Kelly", "ian.kelly@example.com", 28, "+1-312-555-6101", "Address for Ian Kelly", "Company 1", "Designer", "Marketing", "active", "2024-03-01T09:30:00.000", "User note 61", "2024-03-01T10:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["designer","web"],"score":83});
INSERT INTO users VALUES (62, "Jane Sanders", "jane.sanders@example.com", 33, "+1-512-555-6201", "Address for Jane Sanders", "Company 2", "Analyst", "Finance", "active", "2024-03-02T10:30:00.000", "User note 62", "2024-03-02T11:00:00.000", {"theme":"auto","notifications":false,"language":"en"}, {"tags":["analyst"],"score":90});
INSERT INTO users VALUES (63, "Kyle Price", "kyle.price@example.com", 37, "+1-617-555-6301", "Address for Kyle Price", "Company 3", "Engineer", "Engineering", "active", "2024-03-03T11:30:00.000", "User note 63", "2024-03-03T12:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["engineer","elixir"],"score":86});
INSERT INTO users VALUES (64, "Laura Bennett", "laura.bennett@example.com", 30, "+1-305-555-6401", "Address for Laura Bennett", "Company 4", "Manager", "Sales", "active", "2024-03-04T12:30:00.000", "User note 64", "2024-03-04T13:00:00.000", NULL, NULL);
INSERT INTO users VALUES (65, "Mark Wood", "mark.wood@example.com", 42, "+1-619-555-6501", "Address for Mark Wood", "Company 5", "Designer", "Marketing", "active", "2024-03-05T13:30:00.000", "User note 65", "2024-03-05T14:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["designer","veteran"],"score":92});
INSERT INTO users VALUES (66, "Nina Barnes", "nina.barnes@example.com", 25, "+1-212-555-6601", "Address for Nina Barnes", "Company 6", "Analyst", "Finance", "active", "2024-03-06T14:30:00.000", "User note 66", "2024-03-06T15:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["analyst"],"score":81});
INSERT INTO users VALUES (67, "Owen Ross", "owen.ross@example.com", 39, "+1-310-555-6701", "Address for Owen Ross", "Company 7", "Engineer", "Engineering", "active", "2024-03-07T15:30:00.000", "User note 67", "2024-03-07T16:00:00.000", {"theme":"dark","notifications":false,"language":"en","two_factor":true}, {"tags":["engineer","security"],"score":93});
INSERT INTO users VALUES (68, "Paula Henderson", "paula.henderson@example.com", 34, "+1-415-555-6801", "Address for Paula Henderson", "Company 8", "Manager", "Sales", "active", "2024-03-08T16:30:00.000", "User note 68", "2024-03-08T17:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["sales"],"score":84});
INSERT INTO users VALUES (69, "Quincy Coleman", "quincy.coleman@example.com", 27, "+1-206-555-6901", "Address for Quincy Coleman", "Company 9", "Designer", "Marketing", "active", "2024-03-09T17:30:00.000", "User note 69", "2024-03-09T18:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["designer"],"score":80});
INSERT INTO users VALUES (70, "Rachel Jenkins", "rachel.jenkins@example.com", 36, "+1-503-555-7001", "Address for Rachel Jenkins", "Company 0", "Analyst", "Finance", "active", "2024-03-10T18:30:00.000", "User note 70", "2024-03-10T19:00:00.000", {"theme":"auto","notifications":false,"language":"en","email_digest":"monthly"}, {"tags":["analyst","excel"],"score":87});
INSERT INTO users VALUES (71, "Steve Perry", "steve.perry@example.com", 40, "+1-303-555-7101", "Address for Steve Perry", "Company 1", "Engineer", "Engineering", "active", "2024-03-11T09:30:00.000", "User note 71", "2024-03-11T10:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["engineer","senior"],"score":88});
INSERT INTO users VALUES (72, "Tracy Powell", "tracy.powell@example.com", 29, "+1-312-555-7201", "Address for Tracy Powell", "Company 2", "Manager", "Sales", "active", "2024-03-12T10:30:00.000", "User note 72", "2024-03-12T11:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["sales"],"score":82});
INSERT INTO users VALUES (73, "Ursula Long", "ursula.long@example.com", 32, "+1-512-555-7301", "Address for Ursula Long", "Company 3", "Designer", "Marketing", "active", "2024-03-13T11:30:00.000", "User note 73", "2024-03-13T12:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["designer"],"score":85});
INSERT INTO users VALUES (74, "Vincent Patterson", "vincent.patterson@example.com", 38, "+1-617-555-7401", "Address for Vincent Patterson", "Company 4", "Analyst", "Finance", "active", "2024-03-14T12:30:00.000", "User note 74", "2024-03-14T13:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["analyst"],"score":89});
INSERT INTO users VALUES (75, "Wilma Hughes", "wilma.hughes@example.com", 31, "+1-305-555-7501", "Address for Wilma Hughes", "Company 5", "Engineer", "Engineering", "active", "2024-03-15T13:30:00.000", "User note 75", "2024-03-15T14:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["engineer","swift"],"score":86});
INSERT INTO users VALUES (76, "Xander Flores", "xander.flores@example.com", 35, "+1-619-555-7601", "Address for Xander Flores", "Company 6", "Manager", "Sales", "active", "2024-03-16T14:30:00.000", "User note 76", "2024-03-16T15:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["sales","manager"],"score":83});
INSERT INTO users VALUES (77, "Yasmin Washington", "yasmin.washington@example.com", 26, "+1-212-555-7701", "Address for Yasmin Washington", "Company 7", "Designer", "Marketing", "active", "2024-03-17T15:30:00.000", "User note 77", "2024-03-17T16:00:00.000", {"theme":"auto","notifications":false,"language":"en"}, {"tags":["designer"],"score":81});
INSERT INTO users VALUES (78, "Zachary Butler", "zachary.butler@example.com", 43, "+1-310-555-7801", "Address for Zachary Butler", "Company 8", "Analyst", "Finance", "active", "2024-03-18T16:30:00.000", "User note 78", "2024-03-18T17:00:00.000", {"theme":"light","notifications":true,"language":"en","email_digest":"weekly"}, {"tags":["analyst","veteran"],"score":92});
INSERT INTO users VALUES (79, "Abigail Simmons", "abigail.simmons@example.com", 28, "+1-415-555-7901", "Address for Abigail Simmons", "Company 9", "Engineer", "Engineering", "active", "2024-03-19T17:30:00.000", "User note 79", "2024-03-19T18:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["engineer","ruby"],"score":84});
INSERT INTO users VALUES (80, "Bradley Foster", "bradley.foster@example.com", 37, "+1-206-555-8001", "Address for Bradley Foster", "Company 0", "Manager", "Sales", "inactive", "2024-03-20T18:30:00.000", "User note 80", "2024-03-20T19:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["sales"],"score":77});
INSERT INTO users VALUES (81, "Carmen Gonzales", "carmen.gonzales@example.com", 30, "+1-503-555-8101", "Address for Carmen Gonzales", "Company 1", "Designer", "Marketing", "active", "2024-03-21T09:30:00.000", "User note 81", "2024-03-21T10:00:00.000", {"theme":"dark","notifications":true,"language":"es"}, {"tags":["designer","bilingual"],"score":87});
INSERT INTO users VALUES (82, "Derek Bryant", "derek.bryant@example.com", 41, "+1-303-555-8201", "Address for Derek Bryant", "Company 2", "Analyst", "Finance", "active", "2024-03-22T10:30:00.000", "User note 82", "2024-03-22T11:00:00.000", {"theme":"auto","notifications":false,"language":"en"}, {"tags":["analyst"],"score":90});
INSERT INTO users VALUES (83, "Eliza Alexander", "eliza.alexander@example.com", 25, "+1-312-555-8301", "Address for Eliza Alexander", "Company 3", "Engineer", "Engineering", "active", "2024-03-23T11:30:00.000", "User note 83", "2024-03-23T12:00:00.000", NULL, NULL);
INSERT INTO users VALUES (84, "Fernando Russell", "fernando.russell@example.com", 39, "+1-512-555-8401", "Address for Fernando Russell", "Company 4", "Manager", "Sales", "active", "2024-03-24T12:30:00.000", "User note 84", "2024-03-24T13:00:00.000", {"theme":"light","notifications":true,"language":"es"}, {"tags":["sales","bilingual"],"score":85});
INSERT INTO users VALUES (85, "Gabriela Griffin", "gabriela.griffin@example.com", 33, "+1-617-555-8501", "Address for Gabriela Griffin", "Company 5", "Designer", "Marketing", "active", "2024-03-25T13:30:00.000", "User note 85", "2024-03-25T14:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["designer"],"score":86});
INSERT INTO users VALUES (86, "Harold Diaz", "harold.diaz@example.com", 36, "+1-305-555-8601", "Address for Harold Diaz", "Company 6", "Analyst", "Finance", "active", "2024-03-26T14:30:00.000", "User note 86", "2024-03-26T15:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["analyst"],"score":88});
INSERT INTO users VALUES (87, "Isabella Hayes", "isabella.hayes@example.com", 27, "+1-619-555-8701", "Address for Isabella Hayes", "Company 7", "Engineer", "Engineering", "active", "2024-03-27T15:30:00.000", "User note 87", "2024-03-27T16:00:00.000", {"theme":"dark","notifications":true,"language":"en","beta_features":true}, {"tags":["engineer","php"],"score":83});
INSERT INTO users VALUES (88, "Julian Myers", "julian.myers@example.com", 42, "+1-212-555-8801", "Address for Julian Myers", "Company 8", "Manager", "Sales", "active", "2024-03-28T16:30:00.000", "User note 88", "2024-03-28T17:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["sales","veteran"],"score":91});
INSERT INTO users VALUES (89, "Kimberly Ford", "kimberly.ford@example.com", 29, "+1-310-555-8901", "Address for Kimberly Ford", "Company 9", "Designer", "Marketing", "active", "2024-03-29T17:30:00.000", "User note 89", "2024-03-29T18:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["designer","graphics"],"score":84});
INSERT INTO users VALUES (90, "Leonard Hamilton", "leonard.hamilton@example.com", 40, "+1-415-555-9001", "Address for Leonard Hamilton", "Company 0", "Analyst", "Finance", "active", "2024-03-30T18:30:00.000", "User note 90", "2024-03-30T19:00:00.000", {"theme":"auto","notifications":false,"language":"en","two_factor":true}, {"tags":["analyst","security"],"score":93});
INSERT INTO users VALUES (91, "Monica Graham", "monica.graham@example.com", 31, "+1-206-555-9101", "Address for Monica Graham", "Company 1", "Engineer", "Engineering", "active", "2024-03-31T09:30:00.000", "User note 91", "2024-03-31T10:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["engineer","perl"],"score":80});
INSERT INTO users VALUES (92, "Nicholas Sullivan", "nicholas.sullivan@example.com", 34, "+1-503-555-9201", "Address for Nicholas Sullivan", "Company 2", "Manager", "Sales", "active", "2024-04-01T10:30:00.000", "User note 92", "2024-04-01T11:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["sales","manager"],"score":86});
INSERT INTO users VALUES (93, "Ophelia Wallace", "ophelia.wallace@example.com", 28, "+1-303-555-9301", "Address for Ophelia Wallace", "Company 3", "Designer", "Marketing", "active", "2024-04-02T11:30:00.000", "User note 93", "2024-04-02T12:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["designer"],"score":82});
INSERT INTO users VALUES (94, "Patrick Woods", "patrick.woods@example.com", 38, "+1-312-555-9401", "Address for Patrick Woods", "Company 4", "Analyst", "Finance", "active", "2024-04-03T12:30:00.000", "User note 94", "2024-04-03T13:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["analyst"],"score":89});
INSERT INTO users VALUES (95, "Queenie Cole", "queenie.cole@example.com", 26, "+1-512-555-9501", "Address for Queenie Cole", "Company 5", "Engineer", "Engineering", "active", "2024-04-04T13:30:00.000", "User note 95", "2024-04-04T14:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["engineer"],"score":85});
INSERT INTO users VALUES (96, "Raymond West", "raymond.west@example.com", 44, "+1-617-555-9601", "Address for Raymond West", "Company 6", "Manager", "Sales", "active", "2024-04-05T14:30:00.000", "User note 96", "2024-04-05T15:00:00.000", {"theme":"dark","notifications":true,"language":"en","email_digest":"daily"}, {"tags":["sales","senior"],"score":90});
INSERT INTO users VALUES (97, "Sophia Jordan", "sophia.jordan@example.com", 32, "+1-305-555-9701", "Address for Sophia Jordan", "Company 7", "Designer", "Marketing", "active", "2024-04-06T15:30:00.000", "User note 97", "2024-04-06T16:00:00.000", {"theme":"auto","notifications":false,"language":"en"}, {"tags":["designer"],"score":87});
INSERT INTO users VALUES (98, "Theodore Owens", "theodore.owens@example.com", 35, "+1-619-555-9801", "Address for Theodore Owens", "Company 8", "Analyst", "Finance", "active", "2024-04-07T16:30:00.000", "User note 98", "2024-04-07T17:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["analyst"],"score":88});
INSERT INTO users VALUES (99, "Uriel Reynolds", "uriel.reynolds@example.com", 29, "+1-212-555-9901", "Address for Uriel Reynolds", "Company 9", "Engineer", "Engineering", "active", "2024-04-08T17:30:00.000", "User note 99", "2024-04-08T18:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["engineer","c"],"score":84});
INSERT INTO users VALUES (100, "Violet Fisher", "violet.fisher@example.com", 37, "+1-310-555-0011", "Address for Violet Fisher", "Company 0", "Manager", "Sales", "active", "2024-04-09T18:30:00.000", "User note 100", "2024-04-09T19:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["sales","manager"],"score":86});
INSERT INTO users VALUES (101, "Walter Ellis", "walter.ellis@example.com", 41, "+1-415-555-0111", "Address for Walter Ellis", "Company 1", "Designer", "Marketing", "inactive", "2024-04-10T09:30:00.000", "User note 101", "2024-04-10T10:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["designer"],"score":74});
INSERT INTO users VALUES (102, "Xiomara Gibson", "xiomara.gibson@example.com", 27, "+1-206-555-0211", "Address for Xiomara Gibson", "Company 2", "Analyst", "Finance", "active", "2024-04-11T10:30:00.000", "User note 102", "2024-04-11T11:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["analyst"],"score":85});
INSERT INTO users VALUES (103, "Yvonne Hunt", "yvonne.hunt@example.com", 33, "+1-503-555-0311", "Address for Yvonne Hunt", "Company 3", "Engineer", "Engineering", "active", "2024-04-12T11:30:00.000", "User note 103", "2024-04-12T12:00:00.000", {"theme":"auto","notifications":false,"language":"en","compact_view":true}, {"tags":["engineer","haskell"],"score":91});
INSERT INTO users VALUES (104, "Zane Crawford", "zane.crawford@example.com", 39, "+1-303-555-0411", "Address for Zane Crawford", "Company 4", "Manager", "Sales", "active", "2024-04-13T12:30:00.000", "User note 104", "2024-04-13T13:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["sales","manager"],"score":87});
INSERT INTO users VALUES (105, "Adriana Knight", "adriana.knight@example.com", 30, "+1-312-555-0511", "Address for Adriana Knight", "Company 5", "Designer", "Marketing", "active", "2024-04-14T13:30:00.000", "User note 105", "2024-04-14T14:00:00.000", NULL, NULL);
INSERT INTO users VALUES (106, "Benjamin Pierce", "benjamin.pierce@example.com", 36, "+1-512-555-0611", "Address for Benjamin Pierce", "Company 6", "Analyst", "Finance", "active", "2024-04-15T14:30:00.000", "User note 106", "2024-04-15T15:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["analyst"],"score":89});
INSERT INTO users VALUES (107, "Cecilia Berry", "cecilia.berry@example.com", 25, "+1-617-555-0711", "Address for Cecilia Berry", "Company 7", "Engineer", "Engineering", "active", "2024-04-16T15:30:00.000", "User note 107", "2024-04-16T16:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["engineer"],"score":82});
INSERT INTO users VALUES (108, "Dominic Grant", "dominic.grant@example.com", 42, "+1-305-555-0811", "Address for Dominic Grant", "Company 8", "Manager", "Sales", "active", "2024-04-17T16:30:00.000", "User note 108", "2024-04-17T17:00:00.000", {"theme":"light","notifications":true,"language":"en","email_digest":"weekly"}, {"tags":["sales","veteran"],"score":92});
INSERT INTO users VALUES (109, "Esmeralda Wells", "esmeralda.wells@example.com", 28, "+1-619-555-0911", "Address for Esmeralda Wells", "Company 9", "Designer", "Marketing", "active", "2024-04-18T17:30:00.000", "User note 109", "2024-04-18T18:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["designer"],"score":84});
INSERT INTO users VALUES (110, "Franklin Webb", "franklin.webb@example.com", 40, "+1-212-555-1011", "Address for Franklin Webb", "Company 0", "Analyst", "Finance", "active", "2024-04-19T18:30:00.000", "User note 110", "2024-04-19T19:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["analyst"],"score":88});
INSERT INTO users VALUES (111, "Gemma Simpson", "gemma.simpson@example.com", 31, "+1-310-555-1111", "Address for Gemma Simpson", "Company 1", "Engineer", "Engineering", "active", "2024-04-20T09:30:00.000", "User note 111", "2024-04-20T10:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["engineer","dotnet"],"score":86});
INSERT INTO users VALUES (112, "Harrison Stevens", "harrison.stevens@example.com", 34, "+1-415-555-1211", "Address for Harrison Stevens", "Company 2", "Manager", "Sales", "active", "2024-04-21T10:30:00.000", "User note 112", "2024-04-21T11:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["sales","manager"],"score":84});
INSERT INTO users VALUES (113, "Imogen Tucker", "imogen.tucker@example.com", 26, "+1-206-555-1311", "Address for Imogen Tucker", "Company 3", "Designer", "Marketing", "active", "2024-04-22T11:30:00.000", "User note 113", "2024-04-22T12:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["designer"],"score":81});
INSERT INTO users VALUES (114, "Jasper Porter", "jasper.porter@example.com", 38, "+1-503-555-1411", "Address for Jasper Porter", "Company 4", "Analyst", "Finance", "active", "2024-04-23T12:30:00.000", "User note 114", "2024-04-23T13:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["analyst"],"score":91});
INSERT INTO users VALUES (115, "Kendra Hunter", "kendra.hunter@example.com", 29, "+1-303-555-1511", "Address for Kendra Hunter", "Company 5", "Engineer", "Engineering", "active", "2024-04-24T13:30:00.000", "User note 115", "2024-04-24T14:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["engineer"],"score":85});
INSERT INTO users VALUES (116, "Lorenzo Hicks", "lorenzo.hicks@example.com", 43, "+1-312-555-1611", "Address for Lorenzo Hicks", "Company 6", "Manager", "Sales", "active", "2024-04-25T14:30:00.000", "User note 116", "2024-04-25T15:00:00.000", {"theme":"dark","notifications":false,"language":"en","email_digest":"daily"}, {"tags":["sales","veteran"],"score":89});
INSERT INTO users VALUES (117, "Magnolia Crawford", "magnolia.crawford@example.com", 32, "+1-512-555-1711", "Address for Magnolia Crawford", "Company 7", "Designer", "Marketing", "active", "2024-04-26T15:30:00.000", "User note 117", "2024-04-26T16:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["designer"],"score":88});
INSERT INTO users VALUES (118, "Nathaniel Henry", "nathaniel.henry@example.com", 35, "+1-617-555-1811", "Address for Nathaniel Henry", "Company 8", "Analyst", "Finance", "active", "2024-04-27T16:30:00.000", "User note 118", "2024-04-27T17:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["analyst"],"score":87});
INSERT INTO users VALUES (119, "Octavia Boyd", "octavia.boyd@example.com", 27, "+1-305-555-1911", "Address for Octavia Boyd", "Company 9", "Engineer", "Engineering", "active", "2024-04-28T17:30:00.000", "User note 119", "2024-04-28T18:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["engineer","node"],"score":84});
INSERT INTO users VALUES (120, "Preston Mason", "preston.mason@example.com", 41, "+1-619-555-2011", "Address for Preston Mason", "Company 0", "Manager", "Sales", "inactive", "2024-04-29T18:30:00.000", "User note 120", "2024-04-29T19:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["sales"],"score":73});
INSERT INTO users VALUES (121, "Quintessa Morales", "quintessa.morales@example.com", 30, "+1-212-555-2111", "Address for Quintessa Morales", "Company 1", "Designer", "Marketing", "active", "2024-04-30T09:30:00.000", "User note 121", "2024-04-30T10:00:00.000", NULL, NULL);
INSERT INTO users VALUES (122, "Reginald Kennedy", "reginald.kennedy@example.com", 36, "+1-310-555-2211", "Address for Reginald Kennedy", "Company 2", "Analyst", "Finance", "active", "2024-05-01T10:30:00.000", "User note 122", "2024-05-01T11:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["analyst"],"score":90});
INSERT INTO users VALUES (123, "Serena Warren", "serena.warren@example.com", 28, "+1-415-555-2311", "Address for Serena Warren", "Company 3", "Engineer", "Engineering", "active", "2024-05-02T11:30:00.000", "User note 123", "2024-05-02T12:00:00.000", {"theme":"dark","notifications":false,"language":"en","compact_view":true}, {"tags":["engineer","vue"],"score":86});
INSERT INTO users VALUES (124, "Tobias Dixon", "tobias.dixon@example.com", 39, "+1-206-555-2411", "Address for Tobias Dixon", "Company 4", "Manager", "Sales", "active", "2024-05-03T12:30:00.000", "User note 124", "2024-05-03T13:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["sales","manager"],"score":88});
INSERT INTO users VALUES (125, "Unity Marshall", "unity.marshall@example.com", 25, "+1-503-555-2511", "Address for Unity Marshall", "Company 5", "Designer", "Marketing", "active", "2024-05-04T13:30:00.000", "User note 125", "2024-05-04T14:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["designer"],"score":79});
INSERT INTO users VALUES (126, "Valentina Fowler", "valentina.fowler@example.com", 37, "+1-303-555-2611", "Address for Valentina Fowler", "Company 6", "Analyst", "Finance", "active", "2024-05-05T14:30:00.000", "User note 126", "2024-05-05T15:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["analyst"],"score":92});
INSERT INTO users VALUES (127, "Wesley Chambers", "wesley.chambers@example.com", 33, "+1-312-555-2711", "Address for Wesley Chambers", "Company 7", "Engineer", "Engineering", "active", "2024-05-06T15:30:00.000", "User note 127", "2024-05-06T16:00:00.000", {"theme":"light","notifications":true,"language":"en","beta_features":true}, {"tags":["engineer","angular"],"score":85});
INSERT INTO users VALUES (128, "Xyla Rice", "xyla.rice@example.com", 31, "+1-512-555-2811", "Address for Xyla Rice", "Company 8", "Manager", "Sales", "active", "2024-05-07T16:30:00.000", "User note 128", "2024-05-07T17:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["sales"],"score":83});
INSERT INTO users VALUES (129, "Yorick Stone", "yorick.stone@example.com", 40, "+1-617-555-2911", "Address for Yorick Stone", "Company 9", "Designer", "Marketing", "active", "2024-05-08T17:30:00.000", "User note 129", "2024-05-08T18:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["designer","senior"],"score":91});
INSERT INTO users VALUES (130, "Zara Hanson", "zara.hanson@example.com", 26, "+1-305-555-3011", "Address for Zara Hanson", "Company 0", "Analyst", "Finance", "active", "2024-05-09T18:30:00.000", "User note 130", "2024-05-09T19:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["analyst"],"score":80});
INSERT INTO users VALUES (131, "Adrian Ortiz", "adrian.ortiz@example.com", 34, "+1-619-555-3111", "Address for Adrian Ortiz", "Company 1", "Engineer", "Engineering", "active", "2024-05-10T09:30:00.000", "User note 131", "2024-05-10T10:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["engineer","django"],"score":87});
INSERT INTO users VALUES (132, "Brianna Newman", "brianna.newman@example.com", 29, "+1-212-555-3211", "Address for Brianna Newman", "Company 2", "Manager", "Sales", "active", "2024-05-11T10:30:00.000", "User note 132", "2024-05-11T11:00:00.000", {"theme":"auto","notifications":false,"language":"en"}, {"tags":["sales"],"score":82});
INSERT INTO users VALUES (133, "Clarence Garrett", "clarence.garrett@example.com", 42, "+1-310-555-3311", "Address for Clarence Garrett", "Company 3", "Designer", "Marketing", "active", "2024-05-12T11:30:00.000", "User note 133", "2024-05-12T12:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["designer","veteran"],"score":93});
INSERT INTO users VALUES (134, "Delilah Welch", "delilah.welch@example.com", 27, "+1-415-555-3411", "Address for Delilah Welch", "Company 4", "Analyst", "Finance", "active", "2024-05-13T12:30:00.000", "User note 134", "2024-05-13T13:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["analyst"],"score":84});
INSERT INTO users VALUES (135, "Edmund Larson", "edmund.larson@example.com", 38, "+1-206-555-3511", "Address for Edmund Larson", "Company 5", "Engineer", "Engineering", "active", "2024-05-14T13:30:00.000", "User note 135", "2024-05-14T14:00:00.000", {"theme":"auto","notifications":true,"language":"en","two_factor":true}, {"tags":["engineer","security"],"score":94});
INSERT INTO users VALUES (136, "Francesca Frazier", "francesca.frazier@example.com", 32, "+1-503-555-3611", "Address for Francesca Frazier", "Company 6", "Manager", "Sales", "active", "2024-05-15T14:30:00.000", "User note 136", "2024-05-15T15:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["sales","manager"],"score":86});
INSERT INTO users VALUES (137, "Gilbert Burke", "gilbert.burke@example.com", 35, "+1-303-555-3711", "Address for Gilbert Burke", "Company 7", "Designer", "Marketing", "active", "2024-05-16T15:30:00.000", "User note 137", "2024-05-16T16:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["designer"],"score":85});
INSERT INTO users VALUES (138, "Henrietta Hanson", "henrietta.hanson@example.com", 30, "+1-312-555-3811", "Address for Henrietta Hanson", "Company 8", "Analyst", "Finance", "active", "2024-05-17T16:30:00.000", "User note 138", "2024-05-17T17:00:00.000", {"theme":"auto","notifications":false,"language":"en"}, {"tags":["analyst"],"score":89});
INSERT INTO users VALUES (139, "Isaiah Day", "isaiah.day@example.com", 41, "+1-512-555-3911", "Address for Isaiah Day", "Company 9", "Engineer", "Engineering", "inactive", "2024-05-18T17:30:00.000", "User note 139", "2024-05-18T18:00:00.000", NULL, NULL);
INSERT INTO users VALUES (140, "Josephine Sharp", "josephine.sharp@example.com", 28, "+1-617-555-4011", "Address for Josephine Sharp", "Company 0", "Manager", "Sales", "active", "2024-05-19T18:30:00.000", "User note 140", "2024-05-19T19:00:00.000", {"theme":"light","notifications":true,"language":"en"}, {"tags":["sales"],"score":83});
INSERT INTO users VALUES (141, "Kingston Boone", "kingston.boone@example.com", 36, "+1-305-555-4111", "Address for Kingston Boone", "Company 1", "Designer", "Marketing", "active", "2024-05-20T09:30:00.000", "User note 141", "2024-05-20T10:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["designer"],"score":87});
INSERT INTO users VALUES (142, "Lillian Mann", "lillian.mann@example.com", 25, "+1-619-555-4211", "Address for Lillian Mann", "Company 2", "Analyst", "Finance", "active", "2024-05-21T10:30:00.000", "User note 142", "2024-05-21T11:00:00.000", {"theme":"auto","notifications":true,"language":"en"}, {"tags":["analyst"],"score":81});
INSERT INTO users VALUES (143, "Montgomery Mack", "montgomery.mack@example.com", 39, "+1-212-555-4311", "Address for Montgomery Mack", "Company 3", "Engineer", "Engineering", "active", "2024-05-22T11:30:00.000", "User note 143", "2024-05-22T12:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["engineer","rails"],"score":90});
INSERT INTO users VALUES (144, "Nora Williamson", "nora.williamson@example.com", 33, "+1-310-555-4411", "Address for Nora Williamson", "Company 4", "Manager", "Sales", "active", "2024-05-23T12:30:00.000", "User note 144", "2024-05-23T13:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["sales","manager"],"score":85});
INSERT INTO users VALUES (145, "Orlando Stevenson", "orlando.stevenson@example.com", 37, "+1-415-555-4511", "Address for Orlando Stevenson", "Company 5", "Designer", "Marketing", "active", "2024-05-24T13:30:00.000", "User note 145", "2024-05-24T14:00:00.000", {"theme":"auto","notifications":false,"language":"en"}, {"tags":["designer"],"score":88});
INSERT INTO users VALUES (146, "Penelope Ryan", "penelope.ryan@example.com", 31, "+1-206-555-4611", "Address for Penelope Ryan", "Company 6", "Analyst", "Finance", "active", "2024-05-25T14:30:00.000", "User note 146", "2024-05-25T15:00:00.000", {"theme":"light","notifications":true,"language":"en","email_digest":"weekly"}, {"tags":["analyst"],"score":92});
INSERT INTO users VALUES (147, "Quade Spencer", "quade.spencer@example.com", 40, "+1-503-555-4711", "Address for Quade Spencer", "Company 7", "Engineer", "Engineering", "active", "2024-05-26T15:30:00.000", "User note 147", "2024-05-26T16:00:00.000", {"theme":"dark","notifications":false,"language":"en"}, {"tags":["engineer","senior"],"score":91});
INSERT INTO users VALUES (148, "Rosalind Fernandez", "rosalind.fernandez@example.com", 26, "+1-303-555-4811", "Address for Rosalind Fernandez", "Company 8", "Manager", "Sales", "active", "2024-05-27T16:30:00.000", "User note 148", "2024-05-27T17:00:00.000", {"theme":"auto","notifications":true,"language":"es"}, {"tags":["sales","bilingual"],"score":80});
INSERT INTO users VALUES (149, "Sebastian Dunn", "sebastian.dunn@example.com", 34, "+1-312-555-4911", "Address for Sebastian Dunn", "Company 9", "Designer", "Marketing", "active", "2024-05-28T17:30:00.000", "User note 149", "2024-05-28T18:00:00.000", {"theme":"light","notifications":false,"language":"en"}, {"tags":["designer"],"score":86});
INSERT INTO users VALUES (150, "Tabitha Castillo", "tabitha.castillo@example.com", 29, "+1-512-555-5011", "Address for Tabitha Castillo", "Company 0", "Analyst", "Finance", "active", "2024-05-29T18:30:00.000", "User note 150", "2024-05-29T19:00:00.000", {"theme":"dark","notifications":true,"language":"en"}, {"tags":["analyst"],"score":84});

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

-- Test table for large JSON data
CREATE TABLE IF NOT EXISTS big_json_data (
    id INTEGER,
    name STRING,
    small_json JSON,
    medium_json JSON,
    large_json JSON,
    array_json JSON,
    nested_json JSON,
    PRIMARY KEY(id)
);

-- Sample data: Large JSON (5 records)
INSERT INTO big_json_data VALUES (1, "Small JSON test", {"key": "value"}, NULL, NULL, NULL, NULL);
INSERT INTO big_json_data VALUES (2, "Medium JSON test", NULL, {"user": {"name": "Alice", "email": "alice@example.com", "profile": {"age": 30, "city": "San Francisco", "country": "USA", "interests": ["coding", "reading", "hiking"]}}}, NULL, NULL, NULL);
INSERT INTO big_json_data VALUES (3, "Large JSON test", NULL, NULL, {"company": {"name": "TechCorp Inc.", "founded": 2010, "employees": 500, "headquarters": {"address": "123 Tech Street", "city": "San Francisco", "state": "CA", "zip": "94102", "country": "USA"}, "departments": [{"name": "Engineering", "head": "John Smith", "employees": 200, "budget": 5000000}, {"name": "Marketing", "head": "Jane Doe", "employees": 50, "budget": 2000000}, {"name": "Sales", "head": "Bob Wilson", "employees": 100, "budget": 3000000}, {"name": "HR", "head": "Alice Johnson", "employees": 30, "budget": 1000000}], "products": [{"id": "P001", "name": "CloudSync Pro", "price": 99.99, "features": ["Real-time sync", "End-to-end encryption", "Multi-device support", "Offline access", "Team collaboration"]}, {"id": "P002", "name": "DataVault Enterprise", "price": 299.99, "features": ["Unlimited storage", "Advanced analytics", "Custom integrations", "Priority support", "SLA guarantee"]}], "financials": {"revenue": 50000000, "profit": 10000000, "growth_rate": 0.25}}}, NULL, NULL);
INSERT INTO big_json_data VALUES (4, "Array JSON test", NULL, NULL, NULL, {"items": ["item1", "item2", "item3", "item4", "item5", "item6", "item7", "item8", "item9", "item10", "item11", "item12", "item13", "item14", "item15", "item16", "item17", "item18", "item19", "item20", "item21", "item22", "item23", "item24", "item25", "item26", "item27", "item28", "item29", "item30", "item31", "item32", "item33", "item34", "item35", "item36", "item37", "item38", "item39", "item40", "item41", "item42", "item43", "item44", "item45", "item46", "item47", "item48", "item49", "item50"], "numbers": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50]}, NULL);

-- Test table for wide string data
CREATE TABLE IF NOT EXISTS wide_data (
    id INTEGER,
    short_text STRING,
    medium_text STRING,
    long_text STRING,
    very_long_text STRING,
    repeated_text STRING,
    PRIMARY KEY(id)
);

-- Sample data: Wide strings
INSERT INTO wide_data VALUES (1, "Short", "This is a medium length text string for testing purposes.", "This is a longer text string that contains more characters and should test the horizontal scrolling capability of the data pane when displaying wide content in the grid view.", "This is a very long text string that is designed to test the extreme horizontal scrolling capabilities of the application. It contains a lot of text that will definitely exceed the normal display width of most terminal windows. The purpose of this test data is to ensure that users can scroll horizontally to see all the content even when it extends far beyond the visible area. This kind of long text might appear in description fields, log messages, or other free-form text data stored in the database.", "ABCDEFGHIJ ABCDEFGHIJ ABCDEFGHIJ ABCDEFGHIJ ABCDEFGHIJ ABCDEFGHIJ ABCDEFGHIJ ABCDEFGHIJ ABCDEFGHIJ ABCDEFGHIJ");
INSERT INTO wide_data VALUES (2, "Test", "Another medium text.", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam.", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890");
INSERT INTO wide_data VALUES (3, "Data", "Sample text here.", "The quick brown fox jumps over the lazy dog. Pack my box with five dozen liquor jugs. How vexingly quick daft zebras jump!", "The quick brown fox jumps over the lazy dog. Pack my box with five dozen liquor jugs. How vexingly quick daft zebras jump! Sphinx of black quartz, judge my vow. Two driven jocks help fax my big quiz. The five boxing wizards jump quickly. Jackdaws love my big sphinx of quartz. We promptly judged antique ivory buckles for the next prize.", "!@#$%^&*() !@#$%^&*() !@#$%^&*() !@#$%^&*() !@#$%^&*() !@#$%^&*() !@#$%^&*() !@#$%^&*() !@#$%^&*() !@#$%^&*()");

-- Test table for nested JSON data
CREATE TABLE IF NOT EXISTS nested_json (
    id INTEGER,
    name STRING,
    data JSON,
    PRIMARY KEY(id)
);

-- Sample data: Nested JSON structures
INSERT INTO nested_json VALUES (1, "Simple object", {"name": "Alice", "age": 30});
INSERT INTO nested_json VALUES (2, "Nested object", {"user": {"name": "Bob", "profile": {"age": 25, "city": "Tokyo"}}});
INSERT INTO nested_json VALUES (3, "Array of objects", {"users": [{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}, {"id": 3, "name": "Charlie"}]});
INSERT INTO nested_json VALUES (4, "Deep nesting", {"level1": {"level2": {"level3": {"level4": {"level5": {"value": "deep"}}}}}});
INSERT INTO nested_json VALUES (5, "Mixed types", {"string": "hello", "number": 42, "float": 3.14, "boolean": true, "null_val": null, "array": [1, 2, 3], "object": {"key": "value"}});
INSERT INTO nested_json VALUES (6, "Complex structure", {"company": {"name": "TechCorp", "employees": [{"id": 1, "name": "Alice", "department": {"name": "Engineering", "manager": {"name": "Bob", "level": 3}}}, {"id": 2, "name": "Charlie", "department": {"name": "Sales", "manager": {"name": "Diana", "level": 2}}}], "locations": [{"city": "Tokyo", "country": "Japan"}, {"city": "New York", "country": "USA"}]}});
INSERT INTO nested_json VALUES (7, "Long strings in nested", {"article": {"title": "Understanding Database Design Patterns for Modern Applications", "author": {"name": "Dr. Jane Smith", "bio": "Dr. Jane Smith is a renowned database architect with over 20 years of experience in designing scalable systems for Fortune 500 companies. She has authored multiple books on database optimization and speaks regularly at international conferences."}, "content": {"introduction": "In this comprehensive guide, we will explore the fundamental principles of database design that every software engineer should understand. From normalization to denormalization, from ACID compliance to eventual consistency, we cover it all.", "sections": [{"title": "Chapter 1: The Basics", "body": "Database design is both an art and a science. It requires understanding not just the technical aspects but also the business requirements and user behavior patterns that will drive the system."}, {"title": "Chapter 2: Advanced Patterns", "body": "Once you have mastered the basics, it is time to explore advanced patterns such as event sourcing, CQRS, and materialized views. These patterns can significantly improve performance and maintainability of your applications."}], "conclusion": "By following the principles outlined in this guide, you will be well-equipped to design databases that are both performant and maintainable. Remember that the best database design is one that evolves with your application needs."}}});
INSERT INTO nested_json VALUES (8, "Array with long items", {"logs": [{"timestamp": "2024-01-15T10:30:00Z", "level": "ERROR", "message": "Failed to connect to database server at db.example.com:5432. Connection timed out after 30 seconds. Please check network connectivity and firewall rules.", "stacktrace": "at ConnectionPool.getConnection(ConnectionPool.java:142) at DatabaseService.query(DatabaseService.java:89) at UserRepository.findById(UserRepository.java:45)"}, {"timestamp": "2024-01-15T10:30:05Z", "level": "WARN", "message": "Retrying database connection. Attempt 2 of 5. Previous error: Connection refused. The server may be overloaded or temporarily unavailable.", "context": {"retry_count": 2, "max_retries": 5, "backoff_ms": 1000}}, {"timestamp": "2024-01-15T10:30:10Z", "level": "INFO", "message": "Successfully established database connection after 2 retry attempts. Connection pool initialized with 10 connections. Average connection time: 250ms."}]});

-- Test table with very long name (over 40 characters)
CREATE TABLE IF NOT EXISTS this_is_a_very_long_table_name_that_should_definitely_overflow_the_tables_pane_width (
    id INTEGER,
    name STRING,
    PRIMARY KEY(id)
);

INSERT INTO this_is_a_very_long_table_name_that_should_definitely_overflow_the_tables_pane_width VALUES (1, "Test data 1");
INSERT INTO this_is_a_very_long_table_name_that_should_definitely_overflow_the_tables_pane_width VALUES (2, "Test data 2");

-- Child table with very long name
CREATE TABLE IF NOT EXISTS this_is_a_very_long_table_name_that_should_definitely_overflow_the_tables_pane_width.also_a_very_long_child_table_name_for_testing (
    child_id INTEGER,
    value STRING,
    PRIMARY KEY(child_id)
);

INSERT INTO this_is_a_very_long_table_name_that_should_definitely_overflow_the_tables_pane_width.also_a_very_long_child_table_name_for_testing VALUES (1, 1, "Child 1");
INSERT INTO this_is_a_very_long_table_name_that_should_definitely_overflow_the_tables_pane_width.also_a_very_long_child_table_name_for_testing VALUES (1, 2, "Child 2");

-- Index creation
CREATE INDEX IF NOT EXISTS email_idx ON users (email);
CREATE INDEX IF NOT EXISTS name_idx ON users (name);
CREATE INDEX IF NOT EXISTS category_idx ON products (category);
CREATE INDEX IF NOT EXISTS status_idx ON orders (status);
