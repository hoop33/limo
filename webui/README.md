# Search App

A simple search app that uses [ReactiveSearch](https://opensource.appbase.io/reactivesearch/) for UI components and [Elasticsearch](https://www.elastic.co/) as the datastore. The Docker version of the app uses [Cerebro](https://github.com/lmenezes/cerebro) to manage the Elasticseach instance.

## Running the app

#### Without Docker
* [Install and run Elasticseach](https://www.elastic.co/guide/en/elasticsearch/guide/current/running-elasticsearch.html)
* _Optional_: [Install and run Cerebro](https://github.com/lmenezes/cerebro#installation)
* Install the Node modules
```
yarn
```
* Start the app
```
npm start
```

#### With Docker
* [Install and start Docker](https://docs.docker.com/install/)
* Start up the services
```
docker-compose up
```

## Adding data
* Modify the `populateES.sh` script per your requirements
* Run the `populateES.sh` script to populate the ES index