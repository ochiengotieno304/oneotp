## Beba API Reference

#### Register Account

```http
  POST /api/accounts/new
```

| Parameter          | Type     | Description                   |
| :----------------- | :------- | :---------------------------- |
| `name`             | `string` | **Required**. User's name     |
| `email`            | `string` | **Required**. User's email    |
| `password`         | `string` | **Required**. User's password |
| `confirm_password` | `string` | **Required**. User's password |

##### request
```json
  {
    "name": "test name",
    "phone": "+254700001234",
  }
```

##### response
```json
  {
    "name": "test name",
    "phone": "+254700001234",
    "token": "token",
    "secret": "secret"
  }
```

####


