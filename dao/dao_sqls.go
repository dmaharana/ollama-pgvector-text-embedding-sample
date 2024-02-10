package dao

const CreateTextDataStmt = `CREATE TABLE IF NOT EXISTS text_data (
	id BIGSERIAL PRIMARY KEY,
	chunk_id INTEGER NOT NULL,
	file_path TEXT NOT NULL,
	text TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

const InsertTextDataStmt = `INSERT INTO text_data (chunk_id, file_path, text) VALUES %s RETURNING id;`

const InsertMultipleTextDataStmt = `COPY text_data (chunk_id, file_path, text) FROM STDIN WITH CSV HEADER ',' RETURNING id;`

const CreateTextEmbbedingStmt = `CREATE TABLE IF NOT EXISTS text_embbeding (
	id BIGSERIAL PRIMARY KEY,
	text_id INTEGER NOT NULL,
	embbeding VECTOR (%d) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	FOREIGN KEY (text_id) REFERENCES text_data (id)
	);`

const InsertTextEmbbedingStmt = `INSERT INTO text_embbeding (text_id, embbeding) VALUES ($1, $2)  RETURNING id;`

const UpdateTextEmbbedingStmt = `UPDATE text_embbeding SET embbeding = $1 WHERE text_id = $2;`

const SelectNewTextStmt = `SELECT td.id, td.text FROM text_data td
LEFT JOIN text_embbeding te on td.id = te.text_id
WHERE td.id IS NULL`

const SelectNearestTextStmt = `SELECT td.id, td.chunk_id, td.file_path, td.text, te.embbeding FROM text_data td
JOIN text_embbeding te on td.id = te.text_id
WHERE te.text_id in (SELECT text_id FROM text_embbeding ORDER BY embbeding <-> $1 LIMIT $2)`

const SelectNearestTextStmt2 = `SELECT text_id FROM text_embbeding ORDER BY embbeding <-> $1 LIMIT $2`

const SelectCosineSimilarityStmt = `SELECT text_id, 1 - (embbeding <=> $1) as cosine_similarity FROM text_embbeding WHERE (text_id) in (SELECT unnest($2::integer[]))`

// const SelectNewContentStmt = `SELECT td.id FROM text_data td WHERE
// (td.file_path, td.text) NOT IN (SELECT unnest($1::text[])::text), unnest($2::text[])::text)`

const SelectNewContentStmt = `SELECT file_path, text, chunk_id 
FROM 
	unnest($1::text[], $2::text[], $3::integer[]) AS t(file_path, text, chunk_id)
WHERE (file_path, text) NOT IN (SELECT file_path, text FROM text_data)`

const SelectOldContentStmt = `SELECT id, file_path, text, chunk_id 
FROM 
	text_data td
WHERE (file_path, text, chunk_id) IN (SELECT 
	unnest($1::text[]), unnest($2::text[]), unnest($3::integer[]))`
