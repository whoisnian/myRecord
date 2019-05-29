#!/bin/bash
#########################################################################
# File Name: record.sh
# Author: nian
# Blog: https://whoisnian.com
# Mail: zhuchangbao1998@gmail.com
# Created Time: 2019年05月29日 星期三 20时51分34秒
#########################################################################
server_addr="http://127.0.0.1:8000"
access_token="K4X9P2iws28ekGN9"

# record -d new 'a new day-record'
# record -w update 2 'update week-record'
# record -m delete 1
# record -f get 0

op_type="day-record"
method="POST"
id=""
content=""
flag_status=""

if [ $# -gt 1 ]
then
    if [ $1 = "-d" ]
    then
        op_type="day-record"
    elif [ $1 = "-w" ]
    then
        op_type="week-record"
    elif [ $1 = "-m" ]
    then
        op_type="month-record"
    elif [ $1 = "-f" ]
    then
        op_type="flag"
    fi

    if [ $2 = "new" ]
    then
        method="POST"
        content=$3
    elif [ $2 = "update" ]
    then
        method="PUT"
        id=$3
        content=$4
        flag_status=$5
    elif [ $2 = "delete" ]
    then
        method="DELETE"
        id=$3
    elif [ $2 = "get" ]
    then
        method="GET"
        flag_status=$3
    fi
else
    echo "Usage: record -[d|w|m|f] [new|update|delete|get] [id|content]"
    exit
fi

curl -X $method \
  "$server_addr/$op_type/$id?from=0&to=99999999&status=$flag_status" \
  -H "Authorization: Bearer $access_token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "content=$content&status=$flag_status"
