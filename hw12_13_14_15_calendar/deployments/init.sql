CREATE TABLE events (
    id uuid UNIQUE,
    title varchar(256),
    start_time timestamp,
    end_time timestamp,
    description text,
    owner_id int,
    notification_time timestamp null
);
