CREATE TABLE users (
 id SERIAL PRIMARY KEY, -- Идентификатор пользователя
 username VARCHAR(100) NOT NULL, -- Имя пользователя
 email VARCHAR(200) UNIQUE NOT NULL, -- Электронная почта пользователя
 password_hash VARCHAR(255) NOT NULL, -- Хэш пароля пользователя
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время создания кортежа
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Время обновления информации
);

CREATE TABLE roles (
 id SERIAL PRIMARY KEY, -- Идентификатор роли
 name VARCHAR(64) UNIQUE NOT NULL, -- Название роли
 permissions BIGINT NOT NULL, -- Битовая маска возможностей роли
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время создания кортежа
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Время обновления информации
);

CREATE TABLE user_roles (
 user_id INT REFERENCES users(id) ON DELETE CASCADE, -- Идентификатор пользователя
 role_id INT REFERENCES roles(id) ON DELETE CASCADE, -- Идентификатор роли
 PRIMARY KEY (user_id, role_id)
);

CREATE TABLE user_actions (
 id SERIAL PRIMARY KEY, -- Идентификатор действия
 user_id INT REFERENCES users(id) ON DELETE SET NULL, -- Идентификатор пользователя
 action_description VARCHAR(255) NOT NULL, -- Описание действия
 action_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Время действия
);

CREATE INDEX idx_user_actions_user_id on user_actions(user_id);

CREATE TABLE categories (
 id SERIAL PRIMARY KEY, -- Идентификатор категории
 name VARCHAR(100) UNIQUE NOT NULL -- Название категории
);

CREATE TABLE quizzes (
 id SERIAL PRIMARY KEY, -- Идентификатор опроса
 author_id INT REFERENCES users(id) ON DELETE SET NULL, -- Идентификатор автора опроса
 title VARCHAR(255) NOT NULL, -- Название опроса
 description TEXT, -- Описание опроса
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время создания кортежа
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Время обновления информации
);

CREATE INDEX idx_quizzes_author_id on quizzes(author_id);

CREATE TABLE quiz_categories (
 quiz_id INT REFERENCES quizzes(id) ON DELETE CASCADE, -- Идентификатор опроса
 category_id INT REFERENCES categories(id) ON DELETE CASCADE, -- Идентификатор категории
 PRIMARY KEY (quiz_id, category_id)
);

CREATE TABLE questions (
 id SERIAL PRIMARY KEY, -- Идентификатор вопроса
 quiz_id INT REFERENCES quizzes(id) ON DELETE CASCADE, -- Идентификатор опроса
 question_text TEXT NOT NULL, -- Текст вопроса
 question_type VARCHAR(50) CHECK (question_type IN ('choice', 'text')) -- Тип вопроса (выбор ответа/пользовательский ввод)
);

CREATE INDEX idx_questions_quiz_id on questions(quiz_id);

CREATE TABLE text_question_answers (
    question_id INT PRIMARY KEY REFERENCES questions(id) ON DELETE CASCADE, -- Идентификатор вопроса
    right_answer TEXT NOT NULL -- Правильный ответ
);

CREATE TABLE choices (
 id SERIAL PRIMARY KEY, -- Идентификатор варианта ответа
 question_id INT REFERENCES questions(id) ON DELETE CASCADE, -- Идентификатор вопроса
 choice_text VARCHAR(255) NOT NULL, -- Текст варианта ответа
 is_correct BOOLEAN DEFAULT FALSE -- Правильность варианта ответа
);

CREATE INDEX idx_choices_question_id on choices(question_id);

CREATE TABLE choice_question_answers (
    question_id INT PRIMARY KEY REFERENCES questions(id) ON DELETE CASCADE, -- Идентификатор вопроса
    right_choice_id INT REFERENCES choices(id) ON DELETE CASCADE  -- Идентификатор правильного варианта ответа
);

CREATE TABLE choice_answers (
 user_id INT REFERENCES users(id) ON DELETE CASCADE, -- Идентификатор пользователя
 question_id INT REFERENCES questions(id) ON DELETE CASCADE, -- Идентификатор вопроса
 choice_id INT REFERENCES choices(id) ON DELETE CASCADE, -- Идентификатор выбранного варианта ответа (может быть незаданным)
 answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время ответа на вопрос
 PRIMARY KEY (user_id, question_id)
);

CREATE TABLE text_answers (
 user_id INT REFERENCES users(id) ON DELETE CASCADE, -- Идентификатор пользователя
 question_id INT REFERENCES questions(id) ON DELETE CASCADE, -- Идентификатор вопроса
 text_answer TEXT, -- Пользовательский текст ответа (может быть незаданным)
 answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время ответа на вопрос
 PRIMARY KEY (user_id, question_id)
);

CREATE TABLE user_quiz_participations (
 user_id INT REFERENCES users(id) ON DELETE CASCADE, -- Идентификатор пользователя
 quiz_id INT REFERENCES quizzes(id) ON DELETE CASCADE, -- Идентификатор опроса
 participation_count INT DEFAULT 0, -- Количество попыток
 PRIMARY KEY (user_id, quiz_id)
);

CREATE TABLE user_quiz_scores (
 user_id INT REFERENCES users(id) ON DELETE CASCADE, -- Идентификатор пользователя
 quiz_id INT REFERENCES quizzes(id) ON DELETE CASCADE, -- Идентификатор опроса
 score FLOAT, -- Процент выполнения опроса
 last_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Последнее время обновления (необходимо для кэша)
 PRIMARY KEY (user_id, quiz_id)
);

CREATE INDEX idx_user_quiz_scores_quiz_id on user_quiz_scores(quiz_id);

CREATE TABLE quiz_participation_times (
 id SERIAL PRIMARY KEY, -- Идентификатор участия в опросе
 participation_number INT DEFAULT 0, -- Номер попытки
 user_id INT REFERENCES users(id) ON DELETE CASCADE, -- Идентификатор пользователя
 quiz_id INT REFERENCES quizzes(id) ON DELETE CASCADE, -- Идентификатор опроса
 started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время старта участия в опросе
 finished_at TIMESTAMP, -- Время конца участия в опросе
 UNIQUE(user_id, quiz_id, participation_number)
);

CREATE INDEX idx_quiz_participation_times_quiz_id on quiz_participation_times(quiz_id);

CREATE TABLE quiz_statistics (
 quiz_id INT PRIMARY KEY REFERENCES quizzes(id) ON DELETE CASCADE, -- Идентификатор опроса
 total_attempts INT DEFAULT 0, -- Количество попыток
 average_score FLOAT, -- Средний процент выполнения опроса
 average_completion_time INTERVAL, -- Среднее время выполнения опроса
 last_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Последнее время обновления (необходимо для кэша)
);

CREATE TABLE news (
 id SERIAL PRIMARY KEY, -- Идентификатор новости
 author_id INT REFERENCES users(id) ON DELETE CASCADE, -- Идентификатор автора новости
 title VARCHAR(255) NOT NULL, -- Название новости
 news_text TEXT NOT NULL, -- Содержимое новости
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время создания кортежа
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Время обновления информации
);

-- Заполнение тестовыми данными

INSERT INTO users (username, email, password_hash) VALUES
('user1', 'user1@gmail.com', 'hash1'),
('user2', 'user2@gmail.com', 'hash2'),
('user3', 'user3@gmail.com', 'hash3'),
('user4', 'user4@gmail.com', 'hash4'),
('user5', 'user5@gmail.com', 'hash5');

INSERT INTO roles (name, permissions) VALUES
('admin', 3),
('guest', 1);

INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1),
(2, 2),
(3, 2),
(4, 2),
(5, 2);

INSERT INTO user_actions (user_id, action_description) VALUES
(1, 'Logged in'),
(2, 'Logged out'),
(3, 'Changed username'),
(4, 'Created a course'),
(5, 'Complated a course');

INSERT INTO categories (name)
VALUES
  ('Science'),
  ('Math'),
  ('History'),
  ('Literature'),
  ('Art'),
  ('Misc');

INSERT INTO quizzes (author_id, title, description)
VALUES
  (1, 'General Knowledge', 'A quiz on general knowledge topics.'),
  (2, 'Science Quiz', 'Test your science knowledge.'),
  (3, 'Math Quiz', 'Challenging math problems.'),
  (4, 'History Quiz', 'Questions about historical events.'),
  (5, 'Literature Quiz', 'Quiz on literary works.'),
  (1, 'Art Quiz', 'Identify famous artworks.'),
  (2, 'Technology Quiz', 'Latest in tech.'),
  (3, 'Sports Quiz', 'Sports trivia.'),
  (4, 'Music Quiz', 'Questions about music genres.'),
  (5, 'Geography Quiz', 'World geography.');

INSERT INTO quiz_categories (quiz_id, category_id)
VALUES
  (1, 1),
  (2, 1),
  (3, 2),
  (4, 3),
  (5, 4),
  (6, 5),
  (7, 1),
  (8, 6),
  (9, 6),
  (10, 6);

INSERT INTO questions (quiz_id, question_text, question_type)
VALUES
  (1, 'What is the capital of France?', 'choice'),
  (2, 'What is H2O more commonly known as?', 'text'),
  (3, 'Solve: 5 + 7', 'choice'),
  (4, 'Who was the first president of the USA?', 'text'),
  (5, 'Who wrote "1984"?', 'text'),
  (6, 'Which artist painted the Mona Lisa?', 'choice'),
  (7, 'What is the speed of light?', 'text'),
  (8, 'Which country hosted the 2016 Olympics?', 'choice'),
  (9, 'Name a genre of classical music.', 'text'),
  (10, 'What is the largest continent?', 'choice');

INSERT INTO text_question_answers (question_id, right_answer)
VALUES
  (2, 'Water'),
  (4, 'George Washington'),
  (5, 'George Orwell'),
  (7, '299792458 m/s'),
  (9, 'Symphony');

INSERT INTO choices (question_id, choice_text, is_correct)
VALUES
  (1, 'Paris', TRUE),
  (1, 'London', FALSE),
  (3, '12', TRUE),
  (3, '10', FALSE),
  (6, 'Leonardo da Vinci', TRUE),
  (6, 'Vincent van Gogh', FALSE),
  (8, 'Brazil', TRUE),
  (8, 'China', FALSE),
  (10, 'Asia', TRUE),
  (10, 'Africa', FALSE);

INSERT INTO choice_question_answers (question_id, right_choice_id)
VALUES
  (1, 1),
  (3, 3),
  (6, 5),
  (8, 7),
  (10, 9);

INSERT INTO choice_answers (user_id, question_id, choice_id)
VALUES
  (1, 1, 1),
  (2, 3, 3),
  (3, 6, 5),
  (4, 8, 7),
  (5, 10, 9);

INSERT INTO text_answers (user_id, question_id, text_answer)
VALUES
  (1, 2, 'Water'),
  (2, 4, 'George Washington'),
  (3, 5, 'George Orwell'),
  (4, 7, '299792458 m/s'),
  (5, 9, 'Symphony');

INSERT INTO user_quiz_participations (user_id, quiz_id, participation_count)
VALUES
  (1, 1, 1),
  (2, 2, 2),
  (3, 3, 1),
  (4, 4, 3),
  (5, 5, 2);

INSERT INTO user_quiz_scores (user_id, quiz_id, score)
VALUES
  (1, 1, 85.5),
  (2, 2, 90.0),
  (3, 3, 78.0),
  (4, 4, 88.5),
  (5, 5, 92.0);

INSERT INTO quiz_participation_times (participation_number, user_id, quiz_id, started_at, finished_at)
VALUES
  (1, 1, 1, CURRENT_TIMESTAMP - INTERVAL '30 minutes', CURRENT_TIMESTAMP),
  (1, 2, 2, CURRENT_TIMESTAMP - INTERVAL '45 minutes', CURRENT_TIMESTAMP),
  (1, 3, 3, CURRENT_TIMESTAMP - INTERVAL '25 minutes', CURRENT_TIMESTAMP),
  (1, 4, 4, CURRENT_TIMESTAMP - INTERVAL '40 minutes', CURRENT_TIMESTAMP),
  (1, 5, 5, CURRENT_TIMESTAMP - INTERVAL '35 minutes', CURRENT_TIMESTAMP);

INSERT INTO quiz_statistics (quiz_id, total_attempts, average_score, average_completion_time)
VALUES
  (1, 10, 85.0, INTERVAL '30 minutes'),
  (2, 15, 88.5, INTERVAL '40 minutes'),
  (3, 8, 78.0, INTERVAL '25 minutes'),
  (4, 12, 90.0, INTERVAL '35 minutes'),
  (5, 7, 92.0, INTERVAL '45 minutes');

INSERT INTO news (author_id, title, news_text)
VALUES
  (1, 'New Quiz Available!', 'Check out the new quiz on our platform.'),
  (2, 'Weekly Update', 'Here are the latest updates from this week.'),
  (3, 'Maintenance Notice', 'Scheduled maintenance on Saturday night.'),
  (4, 'Feature Release', 'New features have been added to enhance your experience.'),
  (5, 'Community Event', 'Join our upcoming community event.');

-- Procedures
CREATE OR REPLACE PROCEDURE update_user_role(user_id INT, new_role_id INT)
LANGUAGE plpgsql AS $$
BEGIN
  UPDATE user_roles
  SET role_id = new_role_id
  WHERE user_id = user_id;

  INSERT INTO user_actions (user_id, action_description)
  VALUES (user_id, 'Updated role');
END;
$$;

CREATE OR REPLACE PROCEDURE record_quiz_participation(new_user_id INT, new_quiz_id INT)
LANGUAGE plpgsql AS $$
BEGIN
  INSERT INTO quiz_participation_times (user_id, quiz_id, participation_number, started_at)
  VALUES (
    new_user_id,
    new_quiz_id,
    COALESCE((SELECT MAX(participation_number)+1 FROM quiz_participation_times WHERE user_id = new_user_id AND quiz_id = new_quiz_id), 1),
    CURRENT_TIMESTAMP);

  INSERT INTO user_quiz_participations (user_id, quiz_id, participation_count)
  VALUES (new_user_id, new_quiz_id, 1)
  ON CONFLICT (user_id, quiz_id)
  DO UPDATE SET participation_count = user_quiz_participations.participation_count + 1;
END;
$$;

CREATE OR REPLACE PROCEDURE calculate_quiz_statistics(this_quiz_id INT)
LANGUAGE plpgsql AS $$
DECLARE
  this_total_attempts INT;
  this_average_score FLOAT;
  this_average_time INTERVAL;
BEGIN
  SELECT COUNT(*), AVG(finished_at - started_at)
  INTO this_total_attempts, this_average_time
  FROM quiz_participation_times
  WHERE quiz_participation_times.quiz_id = this_quiz_id AND quiz_participation_times.finished_at IS NOT NULL;

  SELECT AVG(score)
  INTO this_average_score
  FROM user_quiz_scores
  WHERE user_quiz_scores.quiz_id = this_quiz_id;

  INSERT INTO quiz_statistics (quiz_id, total_attempts, average_score, average_completion_time, last_update_time)
  VALUES (this_quiz_id, this_total_attempts, this_average_score, this_average_time, CURRENT_TIMESTAMP)
  ON CONFLICT (quiz_id) DO UPDATE
  SET total_attempts = EXCLUDED.total_attempts,
      average_score = EXCLUDED.average_score,
      average_completion_time = EXCLUDED.average_completion_time,
      last_update_time = CURRENT_TIMESTAMP;
END;
$$;

-- Triggers
CREATE OR REPLACE FUNCTION update_user_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_timestamp_trigger
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_user_timestamp();

CREATE OR REPLACE FUNCTION log_quiz_creation()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO user_actions (user_id, action_description)
  VALUES (NEW.author_id, 'Created a quiz: ' || NEW.title);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION log_news_creation()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO user_actions (user_id, action_description)
  VALUES (NEW.author_id, 'New news: ' || NEW.title);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION log_registration()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO user_actions (user_id, action_description)
  VALUES (NEW.id, 'New user: ' || NEW.username);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER log_quiz_creation_trigger
AFTER INSERT ON quizzes
FOR EACH ROW EXECUTE FUNCTION log_quiz_creation();

CREATE TRIGGER log_registration_trigger
AFTER INSERT ON users
FOR EACH ROW EXECUTE FUNCTION log_registration();

CREATE TRIGGER log_news_creation_trigger
AFTER INSERT ON news
FOR EACH ROW EXECUTE FUNCTION log_news_creation();

CREATE OR REPLACE FUNCTION update_quiz_statistics()
RETURNS TRIGGER AS $$
BEGIN
  CALL calculate_quiz_statistics(NEW.quiz_id);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_quiz_statistics_trigger
AFTER INSERT OR UPDATE ON user_quiz_scores
FOR EACH ROW EXECUTE FUNCTION update_quiz_statistics();

CREATE TRIGGER update_quiz_statistics_trigger
AFTER INSERT OR UPDATE ON quiz_participation_times
FOR EACH ROW EXECUTE FUNCTION update_quiz_statistics();