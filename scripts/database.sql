--创建数据库

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

create table store_search
(
    "_id"     serial not null
        constraint store_search_pk
            primary key,
    search    text unique,
    create_at timestamp
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
----

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
    uid_val text;
begin

    select uid from store where title = title_val into uid_val;
    if uid_val is not null
    then
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