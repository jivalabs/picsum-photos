# MySQL database setup

Table structure

Original query for postgresql

```mysql
create table if not exists image 
    (id text primary key, 
     author text not null,
     width integer not null, 
     height integer not null, 
     url text not null);
```

Updated query for mysql
```mysql
create table image 
    (id varchar(30) primary key, 
     author varchar(30) not null,
     width integer not null, 
     height integer not null, 
     url varchar(255) not null);
```