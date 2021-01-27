## jv get

Get the latest jenkins version

### Synopsis

Get the latest jenkins version by querying the maven metadata xml.

```
jv get [flags]
```

### Examples

```
To get the latest weekly release:

    jv get [--username <username> --password <password>]

To get the latest LTS release:

    jv get --version-identifier lts [--username <username> --password <password>]

To get the latest LTS for a particular release train:

    jv get --version-identifier 2.249 [--username <username> --password <password>]

```

### Options

```
  -p, --password string             Password to use
  -u, --url string                  URL to query (default "https://repo.jenkins-ci.org/releases/org/jenkins-ci/main/jenkins-war/maven-metadata.xml")
  -n, --username string             Username to use
  -i, --version-identifier string   The version identifier (default "latest")
```

### Options inherited from parent commands

```
  -v, --debug   Debug Output
      --help    Show help for command
```

### SEE ALSO

* [jv](jv.md)	 - Jenkins Version CLI

