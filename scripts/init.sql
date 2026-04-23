CREATE TABLE IF NOT EXISTS source (
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(100) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO source(name) VALUES ('РИА Новости');


CREATE TABLE IF NOT EXISTS news (
    id              SERIAL PRIMARY KEY,
    title           TEXT NOT NULL,
    link            TEXT UNIQUE NOT NULL,
    content         TEXT,
    source_id       SERIAL NOT NULL,
    status          VARCHAR(64) DEFAULT 'wait',
    likes           INTEGER DEFAULT 0,

    published_at   TIMESTAMP,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT source_fk FOREIGN KEY (source_id) REFERENCES source(id)
);
