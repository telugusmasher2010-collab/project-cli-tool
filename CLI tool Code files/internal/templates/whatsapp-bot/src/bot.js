const { getHandler } = require('./handlers');

/**
 * Dispatches an incoming message to the matching handler plugin.
 * Messages sent by this bot are ignored.
 */
async function handleMessage(sock, db, msg, config) {
  if (msg.key?.fromMe) return;
  const handler = getHandler(msg, config);
  if (!handler) return;
  await handler(sock, db, msg, config);
}

module.exports = { handleMessage };
