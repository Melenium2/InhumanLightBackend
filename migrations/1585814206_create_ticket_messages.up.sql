create table ticket_messages
(
    id bigserial not null PRIMARY KEY,
    who int not null,
    ticket_id int not null,
    message_text VARCHAR not null,
    reply_at TIMESTAMP
);