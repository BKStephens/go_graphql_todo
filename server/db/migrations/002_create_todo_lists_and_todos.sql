CREATE TABLE todo_lists (
  id serial PRIMARY KEY,
  user_id int REFERENCES users ON DELETE CASCADE NOT NULL,
  name text NOT NULL DEFAULT '',
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON todo_lists
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE todos (
  id serial PRIMARY KEY,
  text text NOT NULL DEFAULT '',
  done boolean DEFAULT false,
  todo_list_id int REFERENCES todo_lists ON DELETE CASCADE NOT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON todos
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

---- create above / drop below ----

DROP TABLE todos;
DROP TABLE todo_lists;
