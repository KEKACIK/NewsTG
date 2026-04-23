CREATE TABLE IF NOT EXISTS source (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    link TEXT UNIQUE NOT NULL,
    content TEXT,
    source_id SERIAL NOT NULL,
    CONSTRAINT source_fk FOREIGN KEY (source_id) REFERENCES source(id),
    status VARCHAR(64) DEFAULT 'wait',

    published TIMESTAMP,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE INDEX IF NOT EXISTS idx_news_not_posted ON news(posted) WHERE posted = FALSE;
