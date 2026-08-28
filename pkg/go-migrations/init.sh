#!/bin/sh
if [ -d $GOPATH/tyre-match-backend.com/pkg/go-migrations/template/ ]
then
  cp -a $GOPATH/tyre-match-backend.com/pkg/go-migrations/template/. ./
  exit 1
fi
if [ -d vendor/tyre-match-backend.com/pkg/go-migrations/template ]
then
  cp -a vendor/tyre-match-backend.com/pkg/go-migrations/template/. ./
  exit 1
fi
echo "Dependency path not found"
exit 0