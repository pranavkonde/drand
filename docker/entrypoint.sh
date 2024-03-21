#!bin/sh
## this file is used by the demo projects and CLI for building docker images

set -e
user=drand

if [ -n "$DOCKER_DEBUG" ]; then
   set -x
fi

if [ `id -u` -eq 0 ]; then
    rm -rf ${DRAND_HOME}/.drand
    echo "Changing user to $user"
    # ensure directories are writable
    su-exec "$user" test -w "${DRAND_HOME}" || chown -R -- "$user" "${DRAND_HOME}"
    exec su-exec "$user" "$0" $@
fi

if [ ! -d "${DRAND_HOME}/.drand" -a -n "${DRAND_PUBLIC_ADDRESS}" ]; then
    drand generate-keypair --tls-disable "${DRAND_PUBLIC_ADDRESS}"
fi

exec /usr/local/bin/drand start --verbose --tls-disable --private-listen 0.0.0.0:8080 --control 0.0.0.0:8888
