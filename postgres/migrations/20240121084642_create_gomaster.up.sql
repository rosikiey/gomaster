CREATE TABLE todos (
   id bigserial PRIMARY KEY,
   name TEXT NOT NULL,
   completed BOOLEAN DEFAULT FALSE
);