# API_for_ping_servers

#### Мой первый проект на go, поэтому достаточно простая идея
Написан на простом net/http.<br>
В процессе познакомился с некоторыми основными библиотеками, а также поработал с каналами и конкурентностью
#### Routes

* `GET /api_key`
###### `response`:
```json
{
  "status": "created",
  "api_key": "looooong string"
}
```

* `POST /ping/add`
###### `request json`

```json
{
  "url": "string",
  "email": "string",
  "api_key": "string"
}
```
###### `response`
`Status code 201/400/403/500`

* `DELETE /ping/delete`
###### `request json`

```json
{
  "api_key": "string",
  "url": "string"
}
```
###### `response`
`Status code 200/400/403/500`

* `GET /pings?api_key=loooongstring`
###### `response`

```json
{
  "urls": [
    "string",
    "string"
  ]
}
```