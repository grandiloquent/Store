--创建数据库

create table store_sell
(
    "_id"      serial not null
        constraint store_sell_pk
            primary key,
    uid        text   not null
        constraint store_sell_store_uid_fk
            references store (uid)
            on delete cascade,
    taobao     text,
    wholesaler text,
    quantities int
);



create table store_slide
(
    "_id"     serial not null
        constraint store_slide_pk
            primary key,
    filename  text unique,
    create_at timestamp
);
create table store_category
(
    "_id"     serial not null
        constraint store_category_pk
            primary key,
    name      text unique,
    create_at timestamp
);

--函数

create function store_category_insert(name_val text[]) returns void
    language plpgsql
as
$$
declare
    s text;
begin

    foreach s in array name_val
        loop
            insert into store_category (name, create_at) values (s, now());
        end loop;

end;
$$;
create function store_slide_insert(filename_val text[]) returns void
    language plpgsql
as
$$
declare
    s text;
begin

    foreach s in array filename_val
        loop
            insert into store_slide (filename, create_at) values (s, now());
        end loop;

end;
$$;


---

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

----


create function store_sell_insert(uid_val text,
                                  taobao_val text,
                                  wholesaler_val text,
                                  quantities_val int) returns int
    language plpgsql
as
$$
declare
    id int;
begin
    insert into store_sell (uid, taobao, wholesaler, quantities)
    VALUES (uid_val, taobao_val, wholesaler_val, quantities_val) returning _id into id;
    return id;
end;
$$;

---

CREATE OR REPLACE FUNCTION make_store_uid()
    RETURNS text
    LANGUAGE 'plpgsql'
AS
$BODY$
declare
    new_uid text;
    done    bool;
begin
    done := false;
    while not done
        loop
            select string_agg(x, '')
            into new_uid
            from (
                     select chr(ascii('a') + floor(random() * 26)::integer)
                     from generate_series(1, 6)
                 ) as y(x);
            done := not exists(select 1 from store where uid = new_uid);
        end loop;
    return new_uid;
end;
$BODY$;

---


---

create or replace function store_list(limit_val int, offset_val int)
    returns table
            (
                uid        text,
                title      text,
                price      numeric(15, 6),
                thumbnail  text,
                quantities int
            )
    language plpgsql
as
$$
begin

    return query select store.uid, store.title, store.price, store.thumbnail, ss.quantities
                 from store
                          join store_sell ss on store.uid = ss.uid
                 order by update_at desc
                 limit limit_val
                     offset offset_val;
end;
$$;

---
CREATE OR REPLACE FUNCTION store_update(uid_val text,
                                        title_val text,
                                        price_val numeric,
                                        thumbnail_val text,
                                        details_val text,
                                        specification_val text,
                                        service_val text,
                                        properties_val text[],
                                        showcases_val text[])
    RETURNS void
    LANGUAGE 'plpgsql'
AS
$BODY$
declare
    properties_array text[];
    showcases_array  text[];

begin

    if array_length(properties_val, 1) > 0 then
        properties_array = properties_val;
    end if;

    if array_length(showcases_val, 1) > 0 then
        showcases_array = showcases_val;
    end if;

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
end;
$BODY$;

---

create or replace function store_fetch_details(uid_val text)
    returns table
            (
                title         text,
                price         numeric(15, 6),
                details       text,
                specification text,
                service       text,
                showcases     text[],
                properties    text[],
                taobao        text,
                quantities    int
            )
    language plpgsql
as
$$
begin

    return query select store.title,
                        store.price,
                        store.details,
                        store.specification,
                        store.service,
                        store.showcases,
                        store.properties,
                        ss.taobao,
                        ss.quantities
                 from store
                          left join store_sell ss on store.uid = ss.uid
                 where store.uid = uid_val;
end;
$$;


---

create or replace function store_list(limit_val int, offset_val int)
    returns table
            (
                uid        text,
                title      text,
                price      numeric(15, 6),
                thumbnail  text,
                quantities int
            )
    language plpgsql
as
$$
begin

    return query select store.uid, store.title, store.price, store.thumbnail, ss.quantities
                 from store
                          join store_sell ss on store.uid = ss.uid
                 order by update_at desc
                 limit limit_val
                     offset offset_val;
end;
$$;
---

create or replace function store_list_like_results(like_val text,
                                                   limit_val int,
                                                   offset_val int,
                                                   sort_val int=1)
    returns table
            (
                uid        text,
                title      text,
                price      numeric(15, 6),
                thumbnail  text,
                quantities int
            )
    language plpgsql
as
$$
begin
    if sort_val = 1 then
        return query select s.uid, s.title, s.price, s.thumbnail, ss.quantities
                     from store as s
                              left join store_sell ss on s.uid = ss.uid
                     where s.title like like_val
                        or s.details like like_val
                     order by update_at desc
                     limit limit_val
                         offset offset_val;
    elseif sort_val = 2 then
        return query select s.uid, s.title, s.price, s.thumbnail, ss.quantities
                     from store as s
                              left join store_sell ss on s.uid = ss.uid
                     where s.title like like_val
                        or s.details like like_val
                     order by ss.quantities desc
                     limit limit_val
                         offset offset_val;
    elseif sort_val = 3 then
        return query select s.uid, s.title, s.price, s.thumbnail, ss.quantities
                     from store as s
                              left join store_sell ss on s.uid = ss.uid
                     where s.title like like_val
                        or s.details like like_val
                     order by s.price
                     limit limit_val
                         offset offset_val;
    end if;
end ;
$$;


---------------------

create or replace function store_insert(uid_val text,
                                        title_val text,
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
    if uid_val isnull then
        select uid from store where title = title_val into uid_val;
    end if;
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
