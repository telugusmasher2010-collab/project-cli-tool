const test = require('node:test');
const assert = require('node:assert/strict');
const { getHandler } = require('./handlers');

const MESSAGE = {
  key: { remoteJid: '123456@s.whatsapp.net' },
  message: { conversation: '!ping' },
};

test('ping handler responds with pong', async () => {
  const sent = [];
  const sock = { sendMessage: async (jid, msg) => sent.push({ jid, msg }) };
  const handler = getHandler(MESSAGE, { prefix: '!' });
  assert.ok(handler, 'expected a matching handler for the ping command');
  await handler(sock, null, MESSAGE, { prefix: '!' });
  assert.equal(sent.length, 1);
  assert.equal(sent[0].jid, '123456@s.whatsapp.net');
  assert.equal(sent[0].msg.text, 'pong');
});

test('unknown commands have no handler', () => {
  const msg = { ...MESSAGE, message: { conversation: '!unknown' } };
  assert.equal(getHandler(msg, { prefix: '!' }), null);
});
