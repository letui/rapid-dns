#!/bin/bash
cd web/rapid-ui
npm run build
cp -r dist/* ../ui/
cd ../../
go build

if test -e release;then
  rm -rf release
fi
mkdir release

mv rapid-dns release/