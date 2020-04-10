create table notifications (
    id bigserial not null PRIMARY KEY,
    mess varchar not null,
    created_at TIMESTAMP,
    noti_status VARCHAR not NULL,
    for_user INTEGER not null,
    checked BOOLEAN not null
)