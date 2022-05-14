# Game description

To get dashboard request:

```bash
curl -X GET \
  'http://localhost:8082/api/v1/game/dashboard' \
  --header 'Accept: */*' \
```

To upgrade selected factory, accepted resources: `Iron`, `Copper`, `Gold`. Resources names are case insensitive:

```bash
curl -X POST \
  'http://localhost:8082/api/v1/game/upgrade' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "resource": "Iron"
}'
```
