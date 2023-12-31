-- Создание таблицы пользователей
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       fio VARCHAR,
                       email VARCHAR UNIQUE NOT NULL,
                       password VARCHAR NOT NULL
);

-- Создание таблицы объектов
CREATE TABLE objects (
                         id SERIAL PRIMARY KEY,
                         id_user INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
                         address VARCHAR NOT NULL,
                         is_visible BOOLEAN NOT NULL
);

-- Создание таблицы контрактов
CREATE TABLE contracts (
                           id SERIAL PRIMARY KEY,
                           id_object INT REFERENCES objects(id) ON DELETE CASCADE NOT NULL,
                           data DATE NOT NULL,
                           number VARCHAR NOT NULL,
                           status BOOLEAN NOT NULL
);