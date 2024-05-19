-- Inserting data into the grammars table
INSERT INTO grammars (id, title, created_by, updated_by)
VALUES (1, 'Simple Present', 1, 1),
       (2, 'Present Continuous (Progressive)', 1, 1),
       (3, 'Present Perfect', 1, 1),
       (4, 'Present Perfect Continuous (Progressive)', 1, 1),
       (5, 'Simple Past', 1, 1),
       (6, 'Past Continuous (Progressive)', 1, 1),
       (7, 'Past Perfect', 1, 1),
       (8, 'Past Perfect Continuous (Progressive)', 1, 1),
       (9, 'Simple Future', 1, 1),
       (10, 'Future Continuous (Progressive)', 1, 1),
       (11, 'Future Perfect', 1, 1),
       (12, 'Future Perfect Continuous (Progressive)', 1, 1),
       (13, 'Gerunds', 1, 1),
       (14, 'Infinitives', 1, 1),
       (15, 'Modal Verbs', 1, 1),
       (16, 'Passive Voice', 1, 1),
       (17, 'Reported Speech (Indirect Speech)', 1, 1),
       (18, 'Conditionals (Zero, First, Second, Third, Mixed)', 1, 1),
       (19, 'Articles (Definite and Indefinite)', 1, 1),
       (20, 'Countable and Uncountable Nouns', 1, 1),
       (21, 'Prepositions of Time and Place', 1, 1),
       (22, 'Adjectives and Adverbs', 1, 1),
       (23, 'Comparatives and Superlatives', 1, 1),
       (24, 'Relative Clauses', 1, 1),
       (25, 'Direct and Indirect Objects', 1, 1);

SELECT setval('grammars_id_seq', (SELECT MAX(id) FROM grammars));