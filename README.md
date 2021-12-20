# API Mocker
A tiny JSON rule based API mocker.

The API Mocker intends to simulate an API behavior, supported by a group of rules written as a JSON file structure.

## Dependencies
  [go 1.17+](https://golang.org/doc/install)

## Running
  - Clone this project
  - Prepare your [rule file](#rule-file)
  - from the project folder run `go run . -rules=my_custom_rules.json`

The mock runs on `localhost:8000` by default. These attributes can be changed by using `-host` and `-port` flags.

## Rule file
  The rule file follows a JSON structure with the following minimum fields:
  ```json
  {
    "rules": [
      {
        "endpoint": "/api/auth",
        "method": "GET",
        "items": [
          {
            "response": {
              "status": 200,
              "body": ""
            }
          }
        ]
      }
    ]
  }
  ```

  **where:**
  - `endpoint:` *(mandatory)* an expression path to be included in the url routes. Follows go [HttpRouter Lib schema](https://github.com/julienschmidt/httprouter)
  - `method:` *(mandatory)* any http method (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`).
  - `items:` *(mandatory)* an array of rule items where each item is composed by:
    - `queryString` *(optional)* a regex expression to match against the URL query strings, If not provided, any query string will be accepted.
    - `body:` *(optional)* a regex expression to match against the request body. If not provided, any body will be accepted.
    - `counter`: *(optional)* counter matches the number of calls a request should match (starting from zero). If not provided the matching request will return every time.
    - `response:` *(mandatory)* a structure containing the response to be sent back if the request rules above match, where:
      - `status:` *(mandatory)* a HTTP code
      - `delay:` *(optional)* delay (in ms) to take before sending the response
      - `headers:` *(optional)* an array containing a `"key:value"` list of headers to be sent on the response
      - `body:` *(optional)* a string with the response, with optional dynamic parsable rules based on go [gjson lib](https://github.com/tidwall/gjson)

### Rules by example
  - A GET method where it matches a specific query string
    ```JSON
    {
      "endpoint": "/api/auth",
      "method": "GET",
      "items": [
        {
          "queryString": "foo=1.*bar=2.*|bar=2.*foo=1.*",
          "response": {
            "status": 200
          }
        }
      ]
    }
    ```
  - A GET method simulating an error, after waiting for 5 seconds
    ```JSON
      {
        "endpoint": "/api/auth",
        "method": "GET",
        "items": [
          {
            "response": {
              "status": 413,
              "delay": 5000
            }
          }
        ]
      }
    ```
  - A POST method where the body match specific content
    ```JSON
      {
        "endpoint": "/api/auth",
        "method": "POST",
        "items": [
          {
            "body": "\"name\":",
            "response": {
              "status": 200
            }
          }
        ]
      }
    ```
  - A POST method where the body match multiple content, in this case must have *user* and *password* words
    ```JSON
      {
        "endpoint": "/api/auth",
        "method": "POST",
        "items": [
          {
            "body": ".*user.*password",
            "response": {
              "status": 200
            }
          }
        ]
      }
    ```
  - A POST method simulating error returning with json
    ```JSON
      {
        "endpoint": "/api/auth",
        "method": "POST",
        "items": [
          {
            "response": {
              "status": 417,
              "headers": [
                "Content-Type: application/json"
              ],
              "body": "{\"message\": \"Failed to respond\"}"
            }
          }
        ]
      }
    ```
  - A POST method with dynamic answer, parsing 1st element from request body. Dynamic parsing based on go [gjson lib](https://github.com/tidwall/gjson)
    ```JSON
      {
        "endpoint": "/api/auth",
        "method": "POST",
        "items": [
          {
            "response": {
              "status": 417,
              "headers": [
                "Content-Type: application/json"
              ],
              "body": "{\"requested-hash\": \"{{sha256|@keys|0}}\"}"
            }
          }
        ]
      }
    ```
  - A POST method where the body match specific content with counter. In the first call it returns 200, second matching call 400 and all subsequent call 202 with a json body.
    ```JSON
      {
        "endpoint": "/api/auth",
        "method": "POST",
        "items": [
          {
            "counter": 0,
            "body": ".*user",
            "response": {
              "status": 200
            }
          },
          {
            "counter": 1,
            "body": ".*user",
            "response": {
              "status": 400
            }
          },
          {
            "body": ".*user",
            "response": {
              "status": 202,
              "headers": [
                "Content-Type: application/json"
              ],
              "body": "{\"message\": \"ok\"}"
            }
          }
        ]
      }
    ```
## Docker
Starting the api mocker on docker is very simple:
```
docker build -t api-mocker:latest .
docker run -v "/$(pwd)/rules:/rules" -p 8000:8000 api-mocker:latest  -rules=rules/test-rule.json
```
note that you always need to define the volume and path for the rules.

It is also possible to use it directly from [Docker Hub](https://hub.docker.com/r/nordsec/api-mocker):
`docker pull nordsec/api-mocker`

### Via docker-compose

Example ```docker-compose.yml```:

```yaml
version: "3.7"

services:
  api-mocker:
    build:
      context: .
      target: api-mocker
    command: ["-rules", "rules/test-rule.json"]
    ports:
      - "8000:8000"
    volumes:
      - ./rules:/rules
```

run ```docker-compose up``` wait for it to initialize completely, and visit ```http://localhost:8000/any-path-defined-on-rules-file```
