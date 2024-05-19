-- Insert data
INSERT INTO languages (id, name, code)
VALUES (1, 'English', 'en'),
       (2, 'Spanish', 'es'),
       (3, 'French', 'fr'),
       (4, 'German', 'de'),
       (5, 'Italian', 'it'),
       (6, 'Turkish', 'tr'),
       (7, 'Arabic', 'ar'),
       (8, 'Polish', 'pl'),
       (9, 'Czech', 'cs'),
       (10, 'Hungarian', 'hu'),
       (11, 'Romanian', 'ro'),
       (12, 'Greek', 'el'),
       (13, 'Estonian', 'et'),
       (14, 'Azerbaijani', 'az');

SELECT setval('languages_id_seq', (SELECT MAX(id) FROM languages));