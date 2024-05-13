## Summary

Project is written by using Golang with Echo framework.

Database is PostgreSql and there is ElasticSearch integration and Kibana for monitoring.

There is no security precaution for Kibana. After visiting `localhost:5601` click `Configure Manually` and paste your `ELASTIC_URL`. Kibana will provide a verification code via terminal. After you paste it, you will have access to Kibana.

## Update enviroment file
First update `.env` file to start-up. You can follow the `.env.example`.

Example `.env`:
```
DATABASE="host=db user=postgres password=postgres dbname=insider port=5432"

POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=insider

ELASTIC_URL=http://es:9200

INDEX_NAME=insider
FILE_PATH=sample.json

DEFAULT_PAGE_SIZE=20
```

## Dockerize

To dockerize project run following command;

    docker-compose up    

## Run Locally

Run project with the following command;

    make

Then execute `insider` file with the following command;

    ./insider

## Usage

Project has swagger support. To reach urls visit `localhost:8092/docs/`

You can set order configration by using `configs/ Create Config`

Example Usage: 
```
{
  "is_active": true,
  "sort_option": "click",
  "sort_order": "asc"
}
```

Example ES query:

```
{
  "conditions": [
    {
      "field_name": "click",
      "operation": "lt",
      "value": "500"
    },
    {
      "field_name": "name",
      "operation": "query",
      "value": "ogo"
    }
  ]
}
```