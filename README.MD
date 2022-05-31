

https://user-images.githubusercontent.com/42375242/169675075-c35e6c90-fe35-4142-bc8f-4ff8faacef9e.mp4



## Basic

go mod tidy

make ci-build

bin/darwin_amd64/posto_ipirangad
bin/darwin_amd64/battlemond http --dsn="postgres://postgres:postgres@localhost:32760/posto_ipiranga?sslmode=disable"

docker compose up -d

docker ps -a ('find container name')

docker exec -it posto_ipiranga-postgres bash

psql -h localhost -U postgres -W posto_ipiranga
(Password: 'enter')

posto_ipiranga=# CREATE TABLE....

## NASA

Your API key for cassiano@live.jp is:

rdkAaMXwSX7oBaqLkOlnSjVTJ6nE8lXPm5lEbNiW
You can start using this key to make web service requests. Simply pass your key in the URL when making a web request. Here's an example:

https://api.nasa.gov/planetary/apod?api_key=rdkAaMXwSX7oBaqLkOlnSjVTJ6nE8lXPm5lEbNiW
For additional support, please contact us. When contacting us, please tell us what API you're accessing and provide the following account details so we can quickly find you:

Account Email: cassiano@live.jp
Account ID: dab02599-15ad-43e5-b9fd-2f2ebe6e736c
