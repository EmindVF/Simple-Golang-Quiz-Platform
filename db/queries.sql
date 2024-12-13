-- Пул запросов
SELECT username, email FROM users;

UPDATE users SET password_hash = 'newhash' WHERE username = 'user1';

DELETE FROM users WHERE email LIKE '%@gmail.com';

SELECT name FROM roles WHERE permissions > 1;

SELECT action_description 
FROM user_actions 
WHERE user_id = 1 
ORDER BY action_time DESC 
LIMIT 1;

UPDATE quizzes 
SET description = 'Updated description' 
WHERE author_id = 1;

DELETE FROM categories WHERE name LIKE 'A%';

SELECT title, created_at 
FROM quizzes 
WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '30 days';

SELECT user_id, participation_count 
FROM user_quiz_participations 
WHERE participation_count > 1;

UPDATE user_quiz_scores 
SET score = 95.0 
WHERE user_id = 1 AND quiz_id = 1;

DELETE FROM quiz_participation_times 
WHERE finished_at - started_at > INTERVAL '1 hour';

SELECT title 
FROM news 
WHERE author_id = 2;

SELECT title, created_at 
FROM quizzes 
ORDER BY created_at DESC;

UPDATE roles 
SET permissions = 5 
WHERE name = 'guest';

SELECT question_text 
FROM questions 
WHERE quiz_id = 3;

DELETE FROM user_roles 
WHERE user_id = 5;

SELECT quiz_id, average_score 
FROM quiz_statistics 
WHERE total_attempts > 10;

UPDATE quiz_statistics 
SET average_completion_time = INTERVAL '50 minutes' 
WHERE quiz_id = 1;

DELETE FROM text_answers 
WHERE text_answer = 'Symphony';

SELECT choice_text 
FROM choices 
WHERE question_id = 1;

--

SELECT quiz_id, total_attempts
FROM quiz_statistics
GROUP BY quiz_id, total_attempts
HAVING total_attempts > (SELECT AVG(total_attempts) FROM quiz_statistics);

SELECT user_id, quiz_id
FROM user_quiz_participations
GROUP BY user_id
HAVING COUNT(quiz_id) = 1;

SELECT user_id
FROM user_quiz_scores
WHERE score > ANY(SELECT score FROM user_quiz_scores WHERE quiz_id = 1);

SELECT title, created_at
FROM quizzes
WHERE created_at BETWEEN '2023-01-01' AND '2024-12-31';

SELECT user_id, action_description
FROM user_actions
WHERE action_time BETWEEN CURRENT_TIMESTAMP - INTERVAL '7 days' AND CURRENT_TIMESTAMP;

SELECT category_id
FROM quiz_categories
GROUP BY category_id
HAVING COUNT(DISTINCT quiz_id) = 2;

SELECT question_text
FROM questions
WHERE id BETWEEN 1 AND 5;

SELECT user_id
FROM user_quiz_participations
GROUP BY user_id
HAVING SUM(participation_count) > ANY(SELECT participation_count FROM user_quiz_participations WHERE quiz_id = 1);

SELECT name
FROM roles
WHERE permissions BETWEEN 1 AND 5;

-- Joins

SELECT 
    users.username, 
    roles.name AS role_name
FROM 
    users
INNER JOIN 
    user_roles ON users.id = user_roles.user_id
INNER JOIN 
    roles ON user_roles.role_id = roles.id;

SELECT 
    users.username, 
    roles.name AS role_name
FROM 
    users
LEFT JOIN 
    user_roles ON users.id = user_roles.user_id
LEFT JOIN 
    roles ON user_roles.role_id = roles.id;

SELECT 
    roles.name AS role_name, 
    users.username
FROM 
    roles
RIGHT JOIN 
    user_roles ON roles.id = user_roles.role_id
RIGHT JOIN 
    users ON user_roles.user_id = users.id;

SELECT 
    users.username, 
    roles.name AS role_name
FROM 
    users
FULL OUTER JOIN 
    user_roles ON users.id = user_roles.user_id
FULL OUTER JOIN 
    roles ON user_roles.role_id = roles.id;

SELECT 
    users.username, 
    roles.name AS role_name
FROM 
    users
CROSS JOIN 
    roles;

SELECT 
    u1.username AS user1, 
    u2.username AS user2, 
    substring(u1.email from position('@' in u1.email) + 1) AS email_domain
FROM 
    users u1
JOIN 
    users u2 ON substring(u1.email from position('@' in u1.email) + 1) = substring(u2.email from position('@' in u2.email) + 1)
WHERE 
    u1.id < u2.id;

-- 

SELECT 
    user_id, 
    COUNT(*) AS total_quizzes, 
    AVG(score) AS average_score
FROM 
    user_quiz_scores
GROUP BY 
    user_id;

SELECT 
    user_id, 
    quiz_id, 
    score, 
    RANK() OVER (PARTITION BY quiz_id ORDER BY score DESC) AS rank
FROM 
    user_quiz_scores;

SELECT 
    title
FROM 
    quizzes
UNION
SELECT 
    title
FROM 
    news;

--

SELECT 
    username
FROM 
    users
WHERE 
    EXISTS (
        SELECT 1
        FROM user_quiz_participations
        WHERE user_quiz_participations.user_id = users.id
    );

CREATE TABLE guest_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);
INSERT INTO guest_users (username, email)
SELECT 
    u.username, u.email
FROM 
    users u
JOIN 
    user_roles ur ON u.id = ur.user_id
JOIN 
    roles r ON ur.role_id = r.id
WHERE 
    r.name = 'guest';

SELECT 
    quiz_id, 
    CASE
        WHEN average_score >= 90 THEN 'High'
        WHEN average_score >= 75 THEN 'Medium'
        ELSE 'Low'
    END AS score_category
FROM 
    quiz_statistics;

EXPLAIN ANALYZE
SELECT 
    q.title, u.username AS author
FROM 
    quizzes q
JOIN 
    users u ON q.author_id = u.id;

-- 

SELECT
  u.username,
  COUNT(q.id) AS total_quizzes,
  AVG(us.score) AS average_score,
  MIN(us.score) AS min_score,
  MAX(us.score) AS max_score
FROM
  users u
JOIN
  user_quiz_scores us ON u.id = us.user_id
JOIN
  quizzes q ON us.quiz_id = q.id
GROUP BY
  u.username;

SELECT
  u.username,
  COUNT(qa.quiz_id) AS number_of_attempts
FROM
  users u
JOIN
  user_quiz_participations qa ON u.id = qa.user_id
GROUP BY
  u.username
HAVING
  COUNT(qa.quiz_id) >= 1;

SELECT
  u.username,
  q.title,
  us.score,
  AVG(us.score) OVER (PARTITION BY u.id) AS average_score_for_user
FROM
  users u
JOIN
  user_quiz_scores us ON u.id = us.user_id
JOIN
  quizzes q ON us.quiz_id = q.id;

  --

SELECT 
    users.username, 
    roles.name AS role_name
FROM 
    users
INNER JOIN 
    user_roles ON users.id = user_roles.user_id
INNER JOIN 
    roles ON user_roles.role_id = roles.id;

SELECT
  u.username,
  q.title
FROM
  users u
LEFT JOIN
  quizzes q ON u.id = q.author_id;

SELECT 
    users.username, 
    roles.name AS role_name
FROM 
    users
FULL OUTER JOIN 
    user_roles ON users.id = user_roles.user_id
FULL OUTER JOIN 
    roles ON user_roles.role_id = roles.id;

SELECT 
    users.username, 
    roles.name AS role_name
FROM 
    users
CROSS JOIN 
    roles;

SELECT
  u1.username AS user1,
  u2.username AS user2
FROM
  users u1
JOIN
  users u2 ON u1.id != u2.id;