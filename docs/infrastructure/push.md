---
sidebar_position: 6
---

# Server Push

The `push` module (v0.39) is a WebSocket server-push channel: business code
sends typed messages to users, roles, or everyone, and connected clients
receive them as JSON text frames.

Delivery is **best-effort by contract**: recipients that are offline,
disconnected, or too slow to drain their queue miss the message. Reliable
notification belongs in business storage (a notification table the client
pulls), with the push acting as the real-time hint.

## Sending from Business Code

`push.Notifier` is available from DI. It stays available while the endpoint
is disabled — deliveries are then silently dropped (no connections exist), so
business code never branches on the feature flag.

```go
type OrderService struct {
    notifier push.Notifier
}

func (s *OrderService) NotifyApprovers(ctx context.Context, order *Order, approverIDs []string) error {
    return s.notifier.Push(ctx,
        push.NewMessage("order.pending", map[string]any{"orderId": order.ID}),
        push.ToUsers(approverIDs...),
    )
}
```

| API | Contract |
| --- | --- |
| `Push(ctx, message, targets...)` | delivers the message to every recipient selected by the targets (their union, each connection at most once). A zero message ID or time is filled in. Returns an error only for an invalid message or target set, never for missed recipients |
| `push.NewMessage(type, payload)` | builds a message with a generated ID and the current time |
| `push.ToUsers(userIDs...)` | targets the live connections of the given user IDs |
| `push.ToRoles(roles...)` | targets every connection whose principal holds at least one of the roles (snapshotted at handshake) |
| `push.Broadcast()` | targets every live connection |

Errors: `push.ErrNoTarget` (empty target set or empty users/roles selector —
delivering to nobody is always a caller bug), `push.ErrTypeRequired`
(clients dispatch on the type), `push.ErrUnknownTargetKind` (hand-built
target outside the vocabulary).

### Message envelope

Every push is one JSON text frame:

| Field | Type | Meaning |
| --- | --- | --- |
| `id` | `string` | unique message ID (generated when zero) |
| `type` | `string` | business-defined discriminator clients dispatch on |
| `payload` | any JSON value | arbitrary JSON-serializable payload; omitted when empty |
| `time` | timestamp | send time (filled when zero) |

## The WebSocket Endpoint

The endpoint is opt-in:

```toml
[vef.push]
enabled = true
path = "/ws"                      # default
allowed_origins = []              # empty allows every origin (handshake is token-authenticated anyway)
ping_interval = "30s"             # two missed pongs drop the connection
write_timeout = "10s"             # per outbound frame
send_buffer = 32                  # per-connection outbound queue; too-slow clients are disconnected
max_connections_per_user = 0     # per node; 0 is unlimited
session_recheck_interval = "60s"  # opaque-token revalidation cadence
```

### Client handshake

The handshake is authenticated with the same token mechanism as the API
(`vef.security.token_type`), before any socket exists. Browsers cannot set
an `Authorization` header on the WebSocket API, so the token is accepted
from either channel:

- `Authorization: Bearer <token>` header (non-browser clients), or
- the standard access-token query parameter `__accessToken`.

```js
const ws = new WebSocket(`wss://host/ws?__accessToken=${accessToken}`);
ws.onmessage = (e) => {
    const msg = JSON.parse(e.data); // { id, type, payload, time }
    dispatch(msg.type, msg.payload);
};
```

### Session integration (opaque tokens)

Under `vef.security.token_type = "opaque_token"` the connection is bound to
its login session:

- the session is checked at handshake and re-checked immediately after
  registration, closing the handshake-window race with revocation;
  connections are quarantined from delivery until that recheck passes;
- session revocation (logout, concurrent-login eviction, administrative
  kick) closes the connection immediately through the
  `security.SessionRevocationListener` seam;
- a periodic sweep (`session_recheck_interval`) revalidates every
  connection against the session store, catching expiry.

Terminal close codes are part of the client protocol contract — a client
seeing one must not auto-reconnect (transport-level failures, by contrast,
reconnect with backoff):

| Close code | Constant | Meaning |
| --- | --- | --- |
| `4401` | `push.CloseSessionInvalid` | the login session was revoked or expired; enter the logged-out flow |
| `4429` | `push.CloseTooManyConnections` | the per-user connection cap was reached; do not retry until another connection closes |

Under stateless JWT tokens there is no session to revoke; connections live
until they close or fail heartbeats.

## Multi-Node Relay

On a single node the hub delivers node-locally. In a multi-node deployment,
enabling Redis (`vef.redis.enabled = true`) automatically relays pushes and
revocation kicks across nodes via Redis pub/sub:

- the relay channel is namespaced by Redis database and application name
  (`vef:push:relay:<db>:<app-name>`), so unrelated deployments sharing one
  Redis never cross-deliver;
- an enabled relay refuses to start without `vef.app.name` — the namespace
  would otherwise collide;
- publishes ride a bounded worker, detached from the revocation call path.

No additional configuration is required beyond Redis itself.

## Design Notes

- One goroutine owns each socket's writes; per-connection queues are bounded
  (`send_buffer`) and a client too slow to drain is disconnected rather than
  allowed to exert backpressure on the hub.
- Role targeting matches the roles snapshotted at the connection handshake;
  a role change takes effect on the next connection.
- There is no client-to-server messaging contract: the channel is
  deliberately one-way (clients send only pongs).

## Next Step

Pair pushes with durable state changes through the
[Event Bus](./event-bus) — publish a domain event, store the notification,
then push the real-time hint.
