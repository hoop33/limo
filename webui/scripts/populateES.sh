#! /bin/bash

# Shell script to populate ES Index
# Author: Dinakar Kulkarni
# Github: @dinakar29

## The URL of the ELasticsearch instance
ES_URL=http://localhost:9200

## Directory where JSON date is stored in the form of flat files
JSON_DIR=./sampledata

## Name of the ES Index
INDEX_NAME=searchapp

## The document type, if applicable
DOC_TYPE=document

######################## DO NOT MODIFY BELOW THIS LINE #########################

COUNTER=`date +%s`
for FILE in `find $JSON_DIR/. -type f -name "*.json"`
do {
  curl -H 'Content-Type: application/json' -XPOST $ES_URL/$INDEX_NAME/$DOC_TYPE/$COUNTER -d @$FILE
  COUNTER=$[$COUNTER +1]
}
done

exit 0