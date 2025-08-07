install
```
go install https://github.com/tercnem/jsonserver
```
jalankan 
```
jsonserver -config="./api.json"
```

contoh file config
```
{
  "port": 3000,
  "endpoints": [
    {
      "method": "GET",
      "status": 200,
      "path": "/api",
      "JsonHeader":{},
      "jsonResponse": {"api":"a"}
    },
    {
      "method": "POST",
      "status": 200,
      "path": "/api/2",
      "JsonHeader":{},
      "jsonResponse": {"api":"post"}
    }
  ]
}
```