CREATE TABLE BannersInfo(
    tag INT NOT NULL,
    feature INT NOT NULL,
    is_enabled BOOL NOT NULL,
    b_id INT NOT NULL,
    PRIMARY KEY(b_tag, b_feature),
    FOREIGN KEY (b_id) REFERENCES Banner("id") ON DELETE CASCADE
);
CREATE TABLE Banner(
    "id" SERIAL,
    "value" JSONB NOT NULL,
    PRIMARY KEY ("id")
);