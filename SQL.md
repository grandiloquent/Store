# PostgreSQL Store

## store_search

```
create table store_search
(
    "_id"     serial not null
        constraint store_search_pk
            primary key,
    search    text unique,
    raw       text unique,
    visits    int,
    popular   int,
    create_at timestamp,
    update_at timestamp
);
```

```
create function store_search_insert(search_val text[]) returns void
    language plpgsql
as
$$
declare
    s text;
begin

    foreach s in array search_val
        loop
            insert into store_search (search, create_at) values (s, now());
        end loop;

end;
$$;
```


```
create function store_search_insert(search_val text,
                                    raw_val text,
                                    visits_val int,
                                    popular_val int) returns void
    language plpgsql
as
$$
begin


    insert into store_search (search, raw, visits, popular, create_at, update_at)
    values (search_val, raw_val, visits_val, popular_val, now(), now())
    on conflict (search,raw) do nothing;
end;
$$;
```
