insert into
    users (name, nick, email, password)
values
    (
        "Usuário1",
        "user1",
        "user1@mail.com",
        "$2a$10$finFsyhIR/7UK/8nKmlUu.kdN.Vw3AaHBHBMZlp1HiP3J2JpMgkI6"
    ),
    (
        "Usuário2",
        "user2",
        "user2@mail.com",
        "$2a$10$finFsyhIR/7UK/8nKmlUu.kdN.Vw3AaHBHBMZlp1HiP3J2JpMgkI6"
    ),
    (
        "Usuário3",
        "user3",
        "user3@mail.com",
        "$2a$10$finFsyhIR/7UK/8nKmlUu.kdN.Vw3AaHBHBMZlp1HiP3J2JpMgkI6"
    ),
    (
        "Usuário4",
        "user4",
        "user4@mail.com",
        "$2a$10$finFsyhIR/7UK/8nKmlUu.kdN.Vw3AaHBHBMZlp1HiP3J2JpMgkI6"
    );

insert into
    followers (user_id, follower_id)
values
    (1, 2),
    (3, 1),
    (1, 3);