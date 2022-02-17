CREATE TABLE sessions (
    refreshToken VARCHAR(255) NOT NULL,
    expiresAt integer NOT NULL
);

CREATE TABLE userTable (
    userId serial NOT NULL UNIQUE,
    userName VARCHAR(255) NOT NULL UNIQUE,
    passwordHash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    refreshToken VARCHAR(255),
    expiresAt integer
);

CREATE TABLE deskTable (
    deskId serial NOT NULL UNIQUE ,
    deskName VARCHAR(255) NOT NULL,
    deskDescription VARCHAR(255) NOT NULL
);

CREATE TABLE userDeskCompression (
      userDeskId serial NOT NULL UNIQUE,
      userId integer REFERENCES userTable (userId) ON DELETE CASCADE NOT NULL,
      deskId integer REFERENCES deskTable (deskId) ON DELETE CASCADE NOT NULL
);

CREATE TABLE listTable (
    listId serial NOT NULL UNIQUE,
    listPosition integer NOT NULL,
    listName VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL
);

CREATE TABLE listDeskCompression (
      listDeskId serial NOT NULL UNIQUE,
      deskId integer REFERENCES deskTable (deskId) ON DELETE CASCADE NOT NULL,
      listId integer REFERENCES listTable (listId) ON DELETE CASCADE NOT NULL
);

CREATE TABLE itemTable (
    itemId serial NOT NULL UNIQUE,
    userId integer REFERENCES userTable (userId) ON DELETE CASCADE NOT NULL,
    itemName VARCHAR(255) NOT NULL,
    itemDescription VARCHAR(255) NOT NULL,
    Done BOOLEAN NOT NULL,
    itemPosition integer NOT NULL
);

CREATE TABLE listItemCompression (
      listItemId serial NOT NULL UNIQUE,
      listId integer REFERENCES listTable (listId) ON DELETE CASCADE NOT NULL,
      itemId integer REFERENCES itemTable (itemId) ON DELETE CASCADE NOT NULL
);