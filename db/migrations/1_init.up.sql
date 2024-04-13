CREATE TABLE BannersInfo(
    "id" SERIAL,
    "value" JSONB NOT NULL,
    tag_array INT [] NOT NULL,
    feature INT NOT NULL,
    is_enabled BOOL NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);
CREATE INDEX banner_feature ON BannersInfo USING hash (feature);