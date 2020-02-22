CREATE TABLE tguser (
    id INTEGER PRIMARY KEY,
    state VARCHAR(30),
    select_task VARCHAR(30),
    select_category VARCHAR(30)
);

CREATE TABLE category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30)
);

CREATE TABLE category_tguser (
    category_id INTEGER REFERENCES category(id),
    tguser_id INTEGER REFERENCES tguser(id)
);

CREATE TABLE task (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    description TEXT,
    complete BOOLEAN,
    deadline DATE,
    category_id INTEGER REFERENCES category(id)
); 