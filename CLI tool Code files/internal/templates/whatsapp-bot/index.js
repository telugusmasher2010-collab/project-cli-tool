/**
 * {{project_name}} — WhatsApp Bot Entry Point
 *
 * Initializes the Baileys WhatsApp client, connects to the database,
 * and starts the Fastify HTTP server.
 */

const { makeWASocket, useMultiFileAuthState, DisconnectReason } = require('@whiskeysockets/baileys');
const pino = require('pino');
const qrcode = require('qrcode-terminal');
const Fastify = require('fastify');
const config = require('./config.json');
const { initDatabase } = require('./src/db');
const { handleMessage } = require('./src/bot');

const logger = pino({ level: 'silent' });

async function main() {
  // Initialize database
  const db = initDatabase(config.database.path);
  console.log('[DB] SQLite database ready');

  // Load WhatsApp auth state
  const { state, saveCreds } = await useMultiFileAuthState(config.session.path || './auth_info');

  // Create WhatsApp socket connection
  const sock = makeWASocket({
    printQRInTerminal: true,
    auth: state,
    logger,
    browser: ['{{project_name}}', 'Chrome', '1.0.0'],
  });

  // Handle connection updates
  sock.ev.on('connection.update', (update) => {
    const { connection, lastDisconnect, qr } = update;

    if (qr) {
      console.log('[Bot] Scan the QR code above to connect');
    }

    if (connection === 'close') {
      const statusCode = lastDisconnect?.error?.output?.statusCode;
      const shouldReconnect = statusCode !== DisconnectReason.loggedOut;

      console.log(`[Bot] Connection closed. Status: ${statusCode}`);

      if (shouldReconnect) {
        console.log('[Bot] Reconnecting...');
        main();
      } else {
        console.log('[Bot] Logged out. Delete auth_info/ and restart to reconnect.');
        db.close();
        process.exit(0);
      }
    }

    if (connection === 'open') {
      console.log('[Bot] Connected to WhatsApp');
    }
  });

  // Save credentials on update
  sock.ev.on('creds.update', saveCreds);

  // Handle incoming messages
  sock.ev.on('messages.upsert', async (upsert) => {
    if (upsert.type !== 'notify') return;

    for (const msg of upsert.messages) {
      try {
        await handleMessage(sock, db, msg, config);
      } catch (err) {
        console.error('[Bot] Error handling message:', err.message);
      }
    }
  });

  // Start Fastify HTTP server
  const server = Fastify({ logger: false });

  server.get('/health', async () => ({
    status: 'ok',
    bot: config.botName,
    uptime: process.uptime(),
  }));

  server.get('/stats', async () => {
    const row = db.prepare('SELECT COUNT(*) as count FROM messages').get();
    return { messages: row.count };
  });

  const port = config.server?.port || 3000;
  await server.listen({ port, host: '0.0.0.0' });
  console.log(`[API] Server running on http://localhost:${port}`);

  // Graceful shutdown
  const shutdown = async () => {
    console.log('\n[Bot] Shutting down...');
    await server.close();
    sock.end();
    db.close();
    process.exit(0);
  };

  process.on('SIGINT', shutdown);
  process.on('SIGTERM', shutdown);
}

main().catch((err) => {
  console.error('[Fatal]', err);
  process.exit(1);
});
