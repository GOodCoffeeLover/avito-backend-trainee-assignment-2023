#!/usr/bin/env sh

curl -sL -XPOST localhost:7001/v1/user/11 | jq
curl -sL -XPOST localhost:7001/v1/user/12 | jq
curl -sL -XPOST localhost:7001/v1/user/13 | jq
curl -sL -XPOST localhost:7001/v1/user/14 | jq
curl -sL -XPOST localhost:7001/v1/user/15 | jq
curl -sL -XPOST localhost:7001/v1/segment/SEGMET_1 | jq
curl -sL -XPOST localhost:7001/v1/segment/SEGMET_2 | jq
curl -sL -XPOST localhost:7001/v1/segment/SEGMET_3 | jq
curl -sL -XPOST localhost:7001/v1/segment/SEGMET_4 | jq
curl -sL -XPOST localhost:7001/v1/user/12/assignments -d'{"segments": ["SEGMET_1","SEGMET_2", "SEGMET_3"]}' | jq
curl -sL -XPOST localhost:7001/v1/user/13/assignments -d'{"segments": ["SEGMET_1","SEGMET_2"]}' | jq
curl -sL -XPOST localhost:7001/v1/user/14/assignments -d'{"segments": ["SEGMET_1"]}' | jq

curl -sL -XGET localhost:7001/v1/segment/SEGMET_1 | jq
curl -sL -XGET localhost:7001/v1/user/11 | jq
curl -sL -XGET localhost:7001/v1/user/12/assignments | jq

curl -sL -XGET localhost:7001/v1/segment | jq
curl -sL -XGET localhost:7001/v1/user | jq

curl -sL -XDELETE localhost:7001/v1/segment/SEGMET_2 | jq
curl -sL -XDELETE localhost:7001/v1/user/14 | jq
curl -sL -XDELETE localhost:7001/v1/user/12/assignments -d'{"segments": ["SEGMET_1"]}' | jq


curl -sL -XGET localhost:7001/v1/user/12/assignments | jq
curl -sL -XGET localhost:7001/v1/user/13/assignments | jq
curl -sL -XGET localhost:7001/v1/user/14/assignments | jq

curl -sL -XGET localhost:7001/v1/segment | jq
curl -sL -XGET localhost:7001/v1/user | jq
