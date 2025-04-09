create sequence if not exists seq_tenants start with 1 no cycle;
create table if not exists tenants(
    id bigint primary key default nextval('seq_tenants'),
    uuid char(37) not null,
    code varchar(50) not null,    
    label varchar(100) not null
);
create unique index if not exists idx_tenants_code on tenants(code);
create unique index if not exists idx_tenants_uuid on tenants(uuid);
create index if not exists idx_tenants_label on tenants(label);
create sequence if not exists seq_organizations start with 1 no cycle;
create table if not exists organizations (
    id bigint primary key default nextval('seq_organizations'),
    tenants_id bigint references tenants(id),
    uuid char(37) not null,
    code varchar(50) not null,
    label varchar(50) not null,
    type varchar(10) not null
);
create unique index if not exists idx_organizations_code on organizations(code);
create unique index if not exists idx_organizations_uuid on organizations(uuid);
create index if not exists idx_organizations_label on organizations(label);
create sequence if not exists seq_sectors start with 1 no cycle;
create table if not exists sectors (
    id bigint primary key default nextval('seq_sectors'),
    tenants_id bigint not null not null references tenants(id),
    organizations_id bigint not null not null references organizations(id),
    uuid char(37) not null,
    code varchar(50) not null,
    label varchar(50) not null,
    parent_id bigint references sectors(id),
    has_parent boolean default false,
    depth smallint default 0
);
create index if not exists idx_sectors_code on sectors(code);
create unique index if not exists idx_sectors_uuid on sectors(uuid);
create index if not exists idx_sectors_label on sectors(label);
create sequence if not exists seq_users start with 1 no cycle;
create table if not exists users(
    id bigint primary key default nextval('seq_users'),
    tenants_id bigint not null references tenants(id),
    organizations_id bigint not null references organizations(id),
    uuid char(37) not null,
    last_name varchar(50) not null,
    first_name varchar(50) not null,
    login varchar(50) not null,
    email varchar(50) not null
);
create unique index if not exists idx_users_uuid on users(uuid);
