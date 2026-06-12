kubectl run bombardier-stress \
  --image=alpine/bombardier \
  --namespace=test \
  --restart=Never \
  --overrides='{
    "spec": {
      "containers": [{
        "name": "bombardier-stress",
        "image": "alpine/bombardier",
        "command": [
          "/usr/bin/bombardier",
          "-c", "120",
          "-d", "3m",
          "-H", "Connection: close",
          "-m", "POST",
          "http://cart-api-service:8080/api/v1/carts"
        ],
        "resources": {
          "requests": {"cpu": "100m", "memory": "128Mi"},
          "limits": {"cpu": "200m", "memory": "256Mi"}
        }
      }]
    }
  }'
