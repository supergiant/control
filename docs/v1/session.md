# Session

A Session is a transient record of a [User](user.md) logged-in to the UI.

### Example

#### Request

```json
{
  "user": {
    "username": "username",
    "password": "password"
  }
}
```

#### Response

```json
{
  "id": "generated_session_id_for_cookie",
  "user_id": 1
}
```
