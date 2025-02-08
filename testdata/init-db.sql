CREATE TABLE customers
( id number(10) NOT NULL,
  name varchar2(50) NOT NULL,
  email varchar2(50),
  CONSTRAINT pk PRIMARY KEY (id)
);

INSERT INTO customers(id, name, email) VALUES (2, 'John', 'john@gmail.com');