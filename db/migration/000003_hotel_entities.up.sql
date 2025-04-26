CREATE TABLE pt_hotel (
                          id SERIAL PRIMARY KEY,
                          merchant_id INTEGER DEFAULT 0,
                          mch_id VARCHAR(30) NOT NULL,
                          name VARCHAR(50),
                          type SMALLINT NOT NULL DEFAULT 1,
                          level SMALLINT DEFAULT 0,
                          police_code VARCHAR(15),
                          police_auth_code VARCHAR(60),
                          police_sign VARCHAR(30),
                          province INTEGER DEFAULT 0,
                          city INTEGER DEFAULT 0,
                          area INTEGER DEFAULT 0,
                          address VARCHAR(200),
                          lng VARCHAR(50),
                          lat VARCHAR(50),
                          tel VARCHAR(50),
                          imgurl VARCHAR(255),
                          sales_price NUMERIC(10,2) DEFAULT 0.00,
                          ticket_price NUMERIC(10,2),
                          market_price NUMERIC(10,2),
                          ticket_status SMALLINT DEFAULT 0,
                          order_num INTEGER DEFAULT 0,
                          labels VARCHAR(255),
                          simple_desc VARCHAR(50),
                          open_date VARCHAR(50),
                          decoration_date VARCHAR(50),
                          checkin_start VARCHAR(20),
                          checkin_end VARCHAR(20),
                          room_num INTEGER DEFAULT 0,
                          floor VARCHAR(50),
                          free_park SMALLINT DEFAULT 0,
                          free_park_desc VARCHAR(100),
                          charge_park SMALLINT DEFAULT 0,
                          charge_park_desc VARCHAR(100),
                          child VARCHAR(255),
                          service_person VARCHAR(255),
                          pet VARCHAR(255),
                          network VARCHAR(255),
                          child_facility VARCHAR(255),
                          reception_service VARCHAR(255),
                          cater_service VARCHAR(255),
                          general_facility VARCHAR(255),
                          business_service VARCHAR(255),
                          other_service VARCHAR(255),
                          public_area VARCHAR(255),
                          sport_facility VARCHAR(255),
                          entertainment VARCHAR(255),
                          special_facility VARCHAR(255),
                          weight INTEGER DEFAULT 0,
                          star NUMERIC(3,2) DEFAULT 5.00,
                          status SMALLINT DEFAULT 0,
                          createdAt INTEGER DEFAULT 0,
                          updatedAt INTEGER DEFAULT 0
);

CREATE INDEX idx_hotel_merchant_id ON pt_hotel(merchant_id);

COMMENT ON TABLE pt_hotel IS '酒店';


CREATE TABLE pt_hotel_room_type (
                               id SERIAL PRIMARY KEY,
                               merchant_id INTEGER DEFAULT 0,
                               hotel_id INTEGER DEFAULT 0,
                               name VARCHAR(50),
                               bed VARCHAR(50),
                               sqm VARCHAR(50),
                               floor VARCHAR(50),
                               "window" SMALLINT DEFAULT 0,
                               wifi SMALLINT DEFAULT 0,
                               person VARCHAR(50),
                               smoke SMALLINT DEFAULT 0,
                               facility TEXT,
                               status SMALLINT DEFAULT 0,
                               createdAt INTEGER DEFAULT 0,
                               updatedAt INTEGER DEFAULT 0
);

CREATE INDEX idx_hotel_room_type_hotel_id ON pt_hotel_room_type(hotel_id);
CREATE INDEX idx_hotel_room_type_merchant_id ON pt_hotel_room_type(merchant_id);

COMMENT ON TABLE pt_hotel_room_type IS '酒店客房类型';


CREATE TABLE pt_hotel_room_item (
                                    id SERIAL PRIMARY KEY,
                                    hotel_id INTEGER NOT NULL DEFAULT 0,
                                    merchant_id INTEGER NOT NULL DEFAULT 0,
                                    room_id INTEGER NOT NULL DEFAULT 0,
                                    sn VARCHAR(10) NOT NULL,
                                    build VARCHAR(10),
                                    floor VARCHAR(10) NOT NULL,
                                    status SMALLINT NOT NULL DEFAULT 0,
                                    repair_status SMALLINT DEFAULT 0,
                                    active_status SMALLINT NOT NULL DEFAULT 0,
                                    check_status SMALLINT NOT NULL DEFAULT 1,
                                    clean_status SMALLINT NOT NULL DEFAULT 1,
                                    lock_factory VARCHAR(10) NOT NULL DEFAULT '',
                                    lock_mark VARCHAR(50) NOT NULL DEFAULT '',
                                    lock_device TEXT,
                                    lock_device_status TEXT,
                                    lock_id VARCHAR(30) NOT NULL,
                                    lock_status SMALLINT NOT NULL DEFAULT 0,
                                    lock_online SMALLINT NOT NULL DEFAULT 0,
                                    lock_electric INTEGER NOT NULL DEFAULT 0,
                                    gateway_type INTEGER NOT NULL,
                                    gateway_status SMALLINT NOT NULL DEFAULT 0,
                                    gateway_info TEXT NOT NULL,
                                    createdAt INTEGER,
                                    updatedAt INTEGER
);

CREATE INDEX idx_hotel_room_item_hotel_id ON pt_hotel_room_item(hotel_id);
CREATE INDEX idx_hotel_room_item_merchant_id ON pt_hotel_room_item(merchant_id);
CREATE INDEX idx_hotel_room_item_room_id ON pt_hotel_room_item(room_id);

COMMENT ON TABLE pt_hotel_room_item IS '酒店房间';


CREATE TABLE pt_hotel_sku (
                              id SERIAL PRIMARY KEY,
                              hotel_id INTEGER NOT NULL DEFAULT 0,
                              room_id INTEGER NOT NULL DEFAULT 0,
                              merchant_id INTEGER DEFAULT 0,
                              name VARCHAR(50),
                              shop_price NUMERIC(10,2) DEFAULT 0.00,
                              breakfast_num SMALLINT DEFAULT 0,
                              refund_status SMALLINT DEFAULT 0,
                              refund_audit SMALLINT DEFAULT 0,
                              refund_condition VARCHAR(500) DEFAULT '[]',
                              "desc" TEXT,
                              notice TEXT,
                              sale_num INTEGER DEFAULT 0,
                              status SMALLINT DEFAULT 0,
                              createdAt INTEGER DEFAULT 0,
                              updatedAt INTEGER DEFAULT 0
);

CREATE INDEX idx_hotel_sku_hotel_id ON pt_hotel_sku(hotel_id);
CREATE INDEX idx_hotel_sku_room_id ON pt_hotel_sku(room_id);

COMMENT ON TABLE pt_hotel_sku IS '酒店客房sku';
