drop database if exists dbip;
create database dbip;

\c dbip;

drop table if exists ip_to_city_one;
create table ip_to_city_one (
  ip_start   inet unique,
  ip_end     inet unique,
  continent  char(2),
  country    char(2),
  state_prov text not null default '',
  city       text not null default '',
  latitude   double precision,
  longitude  double precision
);

drop table if exists ip_to_city_two;
create table ip_to_city_two (
  ip_start   inet unique,
  ip_end     inet unique,
  continent  char(2),
  country    char(2),
  state_prov text not null default '',
  city       text not null default '',
  latitude   double precision,
  longitude  double precision
);

drop table if exists config;
create table config (
  last_update timestamp,
  active_table text,
  backup_table text
);

insert into config values(
date_trunc('month', now() - interval '1' month),
'ip_to_city_one',
'ip_to_city_two'
);
