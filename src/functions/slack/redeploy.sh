#!/bin/bash

go build -o /tmp/slack
strip /tmp/slack
zip -j /tmp/pots.zip /tmp/slack
aws lambda update-function-code --function-name pots --zip-file fileb:///tmp/pots.zip
