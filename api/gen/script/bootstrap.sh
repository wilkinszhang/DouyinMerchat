#! /usr/bin/env bash
CURDIR=$(cd $(dirname $0); pwd)
echo "$CURDIR/bin/AuthService"
exec "$CURDIR/bin/AuthService"
