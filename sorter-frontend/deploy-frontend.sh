#!/bin/bash

rm -rf ./deploy

mkdir -p deploy
cp package_deploy.json ./deploy/package.json
cp package-lock.json ./deploy/package-lock.json
cp app.yaml ./deploy/app.yaml

npm install && npm run build
mv ./build ./deploy

cd ./deploy && gcloud app deploy