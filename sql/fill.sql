INSERT INTO stars(name, description, distance, age, magnitude, image, is_active)
VALUES  ('Солнце', 'Наша родная звезда, которая светит нам и греет нас', 0, 5.6, -26.7, 'sun.png', true),
        ('Проксима Центавра', 'Звезда, красный карлик, относящаяся к звёздной системе Альфа Центавра, ближайшая к Солнцу звезда', 4.2, 4.8, 11.1, 'Proxima_Centauri.jpg', true),
        ('Звезда Барнарда', 'Одиночная звезда в созвездии Змееносца', 5.96, 1.5, 9.57, 'Barnard.jpeg', true),
        ('Сириус', 'Ярчайшая звезда ночного неба', 8.6, 3.3, -1.46, 'Sirius.jpg', true),
        ('Лейтен 726-8', 'Двойная звезда в созвездии Кита', 8.73, 5.3, 12.5, 'Leiten.jpg', true);

INSERT INTO users(name, is_moderator)
VALUES  ('Алексей', true),
        ('Сергей', false);

INSERT INTO events(name, status, creation_date, formation_date, completion_date, moderator_id, creator_id)
VALUES ('Превращение в сверхновую', 'fulfilled', '2023-10-14 12:34:56', '2023-10-14 16:34:56', '2023-10-14 18:34:56', 1, 2),
       ('Затухание', 'pending', '2023-10-14 12:34:56', '2023-10-14 16:34:56', '2023-10-14 18:34:56', 1, 2);

INSERT INTO star_events(star_id, event_id)
VALUES (1, 1),
       (3, 2);