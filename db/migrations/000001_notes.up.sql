CREATE TABLE IF NOT EXISTS news (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    title VARCHAR(50) NOT NULL,
    content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS news_categories (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    news_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    CONSTRAINT unique_val UNIQUE (news_id, category_id)
);