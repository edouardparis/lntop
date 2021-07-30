## Docker

To run `lntop` from a docker container:

```sh
# you should first review ./lntop/home/initial-config-template.toml
# note that paths are relevant to situation inside docker and we run under root
# so $HOME directory is /root

# build the container
./build.sh 

# if you have an existing .lntop directory on host machine, you can export it:
# export LNTOP_HOME=~/.lntop

# if you have local lnd node on host machine, point LND_HOME to your actual lnd directory:
export LND_HOME=~/.lnd

# or alternatively if you have remote lnd node, specify paths to auth files explicitly:
# export TLS_CERT_FILE=/path/to/tls.cert
# export MACAROON_FILE=/path/to/readonly.macaroon
# export LND_GRPC_HOST=//<remoteip>:10009

# look into _settings.sh for more details on container configuration

# run lntop from the container
./lntop.sh

# lntop data will be mapped to host folder at ./_volumes/lntop-data
# note that you can review/tweak ./_volumes/lntop-data/config-template.toml after first run
# the ./_volumes/lntop-data/config.toml is the effective (generated) config used by lntop run
```

To see `lntop` logs, you can tail them in another terminal session via:
```sh
./logs.sh -f
```

To start from scratch:
```sh
./clean.sh
./build.sh --no-cache
```