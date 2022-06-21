CREATE TABLE IF NOT EXISTS types(
      id INT NOT NULL PRIMARY KEY,
      type VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS users(
      id UUID PRIMARY KEY,
      first_name VARCHAR(50) NOT NULL,
      last_name VARCHAR(50) NOT NULL,
      login VARCHAR(30) NOT NULL ,
      email TEXT [] NOT NULL,
      bio VARCHAR(50) NOT NULL,
      phone_number TEXT [] NOT NULL,
      type_id INT NOT NULL,
      status BOOLEAN NOT NULL,
      FOREIGN KEY (type_id)
          REFERENCES types(id)
          ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS addresses(
    id uuid NOT NULL PRIMARY KEY,
    user_id uuid,
    country VARCHAR(50) NOT NULL,
    city VARCHAR(50) NOT NULL,
    district VARCHAR(50) NOT NULL ,
    postal_code BIGINT NOT NULL,
    FOREIGN KEY(user_id)
      REFERENCES users(id)
      ON DELETE CASCADE
);