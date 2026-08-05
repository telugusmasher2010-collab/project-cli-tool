/**
 * Simple message-handler registry. Register new commands by appending a
 * handler with a `match` predicate and a `run` function.
 */
const handlers = [
  {
    match: (msg, config) => {
      const text = msg.message?.conversation || '';
      return text.startsWith((config.prefix || '!') + 'ping');
    },
    run: async (sock, _db, msg) => {
      await sock.sendMessage(msg.key.remoteJid, { text: 'pong' });
    },
  },
];

/**
 * Returns the first handler whose predicate matches the incoming message,
 * or null when no handler applies.
 */
function getHandler(msg, config) {
  const handler = handlers.find((h) => h.match(msg, config));
  return handler ? handler.run : null;
}

module.exports = { getHandler };
