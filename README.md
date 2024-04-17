# Reading List with Go


#### Pull postgres docker image: 
```bash
 docker pull postgres
 ```

 #### Confirm installation
 ```bash
 docker images | grep postgres
 ```

 #### Setup and run container (Update Password you want to use)
 ```bash
 docker run --name reading-list-db -e POSTGRES_PASSWORD=secure -d -p 5432:5432 postgres
```

#### Login to newly created DB
```bash
psql -h localhost -p 5432 -U postgres
```
and enter password you entered in previous step

#### Create DB for reading list
```bash
CREATE DATABASE readinglist;
```

#### Create role to access db (to not to give default db password)
```bash
CREATE ROLE readinglist WITH LOGIN PASSWORD 'password';
```

#### Switch to newly created db from default postgres db
```bash
\c readinglist;
```

#### Create table with following
```bash
CREATE TABLE IF NOT EXISTS books(
    id bigserial PRIMARY KEY, 
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL, 
    published integer NOT NULL, 
    pages integer NOT NULL, 
    genres text[] NOT NULL, 
    rating real NOT NULL,
    version integer NOT NULL DEFAULT 1
);
```

#### Finally, grant `books` CRUD access to `readinglist` role and sequence access.
```bash
GRANT SELECT, INSERT, UPDATE, DELETE ON books TO readinglist;
```

```bash
GRANT USAGE,SELECT ON SEQUENCE books_id_seq TO readinglist;
```
