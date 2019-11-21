# PostgreSQL Store

## store

```
drop table store CASCADE;

create table store
(
    "_id"         serial         not null
        constraint store_pk
            primary key,
    uid           text           not null,
    title         text           not null,
    price         numeric(15, 6) not null,
    thumbnail     text,
    details       text,
    specification text,
    service       text,
    properties    text[],
    showcases     text[],
    create_at     timestamp,
    update_at     timestamp
);

```

```

create or replace function store_insert(title_val text,
                                        price_val numeric,
                                        thumbnail_val text,
                                        details_val text,
                                        specification_val text,
                                        service_val text,
                                        properties_val text[],
                                        showcases_val text[]) returns text
    language plpgsql
as
$$
declare
    uid_val          text;
    properties_array text[];
    showcases_array  text[];

begin

    if array_length(properties_val, 1) > 0 then
        properties_array = properties_val;
    end if;

    if array_length(showcases_val, 1) > 0 then
        showcases_array = showcases_val;
    end if;

    select uid from store where title = title_val into uid_val;
    if uid_val is not null then
        update
            store
        set title= COALESCE(title_val, title),
            price = COALESCE(price_val, price),
            thumbnail = COALESCE(thumbnail_val, thumbnail),
            details = COALESCE(details_val, details),
            specification = COALESCE(specification_val, specification),
            service = COALESCE(service_val, service),
            properties = COALESCE(properties_array, properties),
            showcases = COALESCE(showcases_array, showcases),
            update_at = now()
        where uid = uid_val;
        return uid_val;
    end if;
    insert into store(uid,
                      title,
                      price,
                      thumbnail,
                      details,
                      specification,
                      service,
                      properties,
                      showcases,
                      create_at,
                      update_at)
    values (make_store_uid(),
            title_val,
            price_val,
            thumbnail_val,
            details_val,
            specification_val,
            service_val,
            properties_val,
            showcases_val,
            now(),
            now()) returning uid into uid_val;
    return uid_val;
end;
$$;
```

## store_search

```
drop table store_search;

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
create or replace function store_search_insert(search_val text,
                                    raw_val text,
                                    visits_val int,
                                    popular_val int) returns void
    language plpgsql
as
$$
begin


    insert into store_search (search, raw, visits, popular, create_at, update_at)
    values (search_val, raw_val, visits_val, popular_val, now(), now())
    on conflict (search) do nothing;
end;
$$;
 
- SELECT * FROM store_search_insert('U盘','',0,0)
```

```
create or replace function store_fetch_keywords(limit_val int,
                                                offset_val int,
                                                sort_val int=1)
    returns table
            (
                search text
            )
    language plpgsql
as
$$
begin
    if sort_val = 1 then
        return query select store_search.search
                     from store_search
                     order by update_at desc
                     limit limit_val
                         offset offset_val;
    elseif sort_val = 2 then
        return query select store_search.search
                     from store_search
                     order by popular desc
                     limit limit_val
                         offset offset_val;
    elseif sort_val = 3 then
        return query select store_search.search
                     from store_search
                     order by visits desc
                     limit limit_val
                         offset offset_val;
    end if;
end ;
$$;

select *
from store_fetch_keywords(6, 0);
```

