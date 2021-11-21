#! /bin/bash

set -e

SEPARATOR="
--------------------------------------------------------------------------------
"

echo "${SEPARATOR}"
echo "
Building posts
"
if ! make build-posts; then
	echo "errored posts' building"
	exit 1
fi

echo "${SEPARATOR}"
echo "
Building users
"
if ! make build-users; then
	echo "errored users' building"
	exit 1
fi

echo "${SEPARATOR}"
echo "
Building gateway
"
if ! make build-gateway-release; then
	echo "errored users' building"
	exit 1
fi

echo ${SEPARATOR}
echo "
Bundling
"

BUNDLE="bundle"

mkdir -p "${BUNDLE}"
cp prod.nix "${BUNDLE}/default.nix"
cp .env "${BUNDLE}"
cp docker-compose-infra.yml "${BUNDLE}"
cp -r nginx/ "${BUNDLE}"
cp post/cmd/posts/posts "${BUNDLE}"
cp user/cmd/users/users "${BUNDLE}"
cp gateway/target/release/gateway "${BUNDLE}"
cp -r front "${BUNDLE}"
rsync -avz -e "ssh -i ${PEM_FILE}" "${BUNDLE}" "${SERVER_USER}@${SERVER_IP}:/home/${SERVER_USER}/blog"
ssh "${SERVER_USER}@${SERVER_IP}" -i "${PEM_FILE}" "cd ~/blog/bundle && echo 'use nix' > .envrc && direnv allow"

rm -rf "${BUNDLE}"
