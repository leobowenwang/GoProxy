# create table

```postgresql
CREATE TABLE client
(
    username TEXT  NOT NULL UNIQUE PRIMARY KEY,
    password TEXT  NOT NULL,
    roles    JSONB NOT NULL
);
```

# insert data

```postgresql
INSERT INTO public.client (username, password, roles)
VALUES ('leo', 'rawr', '{"admin","lion"}');

INSERT INTO public.client (username, password, roles)
VALUES ('joe', '$2a$12$GX6IM71bWwM.CHsWjXjnI.TzyEjw2IrAJgm26TSbuwjEWaBm0F5Fe
', '{"not-admin"}');


```