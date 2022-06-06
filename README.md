# ecsh
ecsh is a tool to execute ECS Exec with ease.

# Requirement
- [session-manager-plugin](https://docs.aws.amazon.com/ja_jp/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)

# Install
## Mac
```
$ ECSH_VERSION=0.0.3
$ curl -OL https://github.com/kishii4726/ecsh/releases/download/v${ECSH_VERSION}/ecsh_v${ECSH_VERSION}_darwin_amd64.zip

$ unzip ecsh_v${ECSH_VERSION}_darwin_amd64.zip ecsh

$ sudo cp ecsh /usr/local/bin
```

## Linux
```
$ ECSH_VERSION=0.0.3
$ curl -OL https://github.com/kishii4726/ecsh/releases/download/v${ECSH_VERSION}/ecsh_v${ECSH_VERSION}_linux_amd64.zip

$ unzip ecsh_v${ECSH_VERSION}_linux_amd64.zip ecsh

$ sudo cp ecsh /usr/local/bin
```

# Usage
```
$ ecsh
```

![ecsh_v0 0 3](https://user-images.githubusercontent.com/46281949/172080245-6cbf0a2e-74aa-49fe-ae81-811b0485a1c0.gif)
