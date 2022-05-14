# Game description

To get dashboard request:

```bash
curl -X GET \
  'http://localhost:8082/api/v1/game/dashboard' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)'
```

To upgrade selected factory, accepted resources: `Iron`, `Copper`, `Gold`. Resources names are case insensitive:

```bash
curl -X POST \
  'http://localhost:8082/api/v1/game/upgrade' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "resource": "Iron"
}'
```
