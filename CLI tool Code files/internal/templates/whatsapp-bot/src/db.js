const fs = require('fs');
const path = require('path');
const Database = require('better-sqlite3');

/**
 * Opens (or creates) the SQLite database at dbPath and applies the schema
 * shipped in database/schema.sql.
 */
function initDatabase(dbPath) {
  fs.mkdirSync(path.dirname(dbPath), { recursive: true });
  const db = new Database(dbPath);
  const schema = fs.readFileSync(
    path.join(__dirname, '..', 'database', 'schema.sql'),
    'utf8',
  );
  db.exec(schema);
  return db;
}

module.exports = { initDatabase };
