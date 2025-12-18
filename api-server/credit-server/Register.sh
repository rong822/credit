#!/bin/bash
#
USER=$(date '+T1%Y%m%d%H%M%S')
echo "POST request Regist"
echo
RESPONSE=$(curl -s -X POST \
  http://localhost:8787/api/register \
  -H "content-type: application/x-www-form-urlencoded" \
  -d "username=$USER&password=testpassword")
echo $RESPONSE
RESPONSE=$(echo $RESPONSE | jq ".token" | sed "s/\"//g")
echo
echo "ORG1 token is $RESPONSE"
echo