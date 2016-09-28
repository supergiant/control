# User

A User is pretty straightforward. The `role` can be "admin" or "user", the
latter being restricted from creating additional Users.

### Example

#### Request

```json
{
  "username": "username",
  "password": "password",
  "role": "user"
}
```

#### Response

```json
{
  "username": "username",
  "role": "user",
  "api_token": "GENERATED_API_TOKEN"
}
```
