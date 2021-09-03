# HTTP Request

1 - Create the file `example.yml` by runnig `$ httpreq new example`. You'll get
the following content:

```
given:
  https: false
  host: localhost
  port: "8080"
  headers:
    Authorization: Bearer my.awesome.token
then:
- call: POST /api/v1/login
  id: login
  with_headers:
    X-Parameter: magic
  and_send: '{ "email": "mr.frodo@lotr.com", "password": "{{ .password }}" }'
where:
  email: mr.frodo@lotr.com
  password: $PASSWORD
```

2 - On another terminal start the echo server with `$ httpreq echo`;

2 - Then run `$ PASSWORD=123456 httpreq exec example.yml`;
