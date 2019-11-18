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

--函数

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
