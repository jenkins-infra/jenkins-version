## jv download

Download the latest jenkins version

### Synopsis

Download the latest jenkins version by querying the maven metadata xml.

```
jv download [flags]
```

### Examples

```
To download the latest weekly release:

    jv download [--username <username> --password <password>]

To download the latest LTS release:

    jv download --version-identifier lts [--username <username> --password <password>]

To download the latest LTS for a particular release train:

    jv download --version-identifier 2.249 [--username <username> --password <password>]

```

### Options

```
  -p, --password string             Password to use (envVar: MAVEN_REPOSITORY_PASSWORD)
  -u, --url string                  URL to query (envVar: JENKINS_DOWNLOAD_URL) (default "https://repo.jenkins-ci.org/releases/org/jenkins-ci/main/jenkins-war/")
  -n, --username string             Username to use (envVar: MAVEN_REPOSITORY_USERNAME)
  -i, --version-identifier string   The version identifier (envVar: JENKINS_VERSION) (default "latest")
  -w, --war string                  The location to download the war file to (envVar: WAR) (default "/tmp/jenkins.war")
```

### Options inherited from parent commands

```
  -v, --debug   Debug Output
      --help    Show help for command
```

### SEE ALSO

* [jv](jv.md)	 - Jenkins Version CLI

