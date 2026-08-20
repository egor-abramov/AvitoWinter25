CREATE TABLE IF NOT EXISTS users
(
    id              uuid PRIMARY KEY                       DEFAULT gen_random_uuid(),
    username        character varying(255) UNIQUE NOT NULL,
    hashed_password character varying(60)         NOT NULL,
    coins           integer                       NOT NULL default 1000 CHECK ( coins >= 0 )
);


CREATE TABLE IF NOT EXISTS transaction
(
    id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user uuid REFERENCES users (id),
    to_user   uuid REFERENCES users (id),
    amount    integer not null CHECK ( amount > 0 )
);

CREATE INDEX idx_from ON transaction (from_user);
CREATE INDEX idx_to ON transaction (to_user);

CREATE TABLE merch
(
    id    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name  character varying(31) UNIQUE NOT NULL,
    price integer                      NOT NULL CHECK ( price > 0 )
);

CREATE TABLE IF NOT EXISTS user_merch
(
    id       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id  uuid REFERENCES users (id),
    merch_id uuid REFERENCES merch (id)
);