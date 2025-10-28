# Introduction

This utility helps automate A-B-C testing for analysis plugin migration. It compares three different versions:

- **Version A** - Pure analysis plugin without any modifications
- **Version B** - Analysis plugin with disabled JSAnalysis checks
- **Version C** - ESLint plugin

## Preparation Steps

### Analysis Plugin Preparation

To use the analysis plugin in this testing utility, you must have separate versions for each variant.

To build a new version of the analysis plugin:
1. Update the version in `pom.xml`
2. Update the version in `maven-plugin/pom.xml`
3. Build with `mvn clean install`

**Example:**

For root `pom.xml`:
```xml
<groupId>com.sap.ca</groupId>
<artifactId>analysis-plugin-parent</artifactId>
<version>2.1.8-A</version>
<packaging>pom</packaging>
```

For `maven-plugin/pom.xml`:
```xml
<parent>
	<groupId>com.sap.ca</groupId>
	<artifactId>analysis-plugin-parent</artifactId>
	<version>2.1.8-A</version>
	<relativePath>../pom.xml</relativePath>
</parent>
```

After building the new versions, verify that your `.m2` repository actually contains these versions:
```sh
ls ~/.m2/repository/com/sap/ca/analysis-plugin
```

Note: Your `.m2` repository location may be different.

### Analysis Plugin Standalone Preparation (Optional)

The standalone version is a modified version of the analysis plugin that can test `JSAnalysis` checks against a single file.

**!!WARNING!!**
For now standalone checks run only for `JSAnalysis` checks, so to disable `JSAnalysis` checks you actually should to comment all rules for `JSAnalysis` instantiating and then build your standalone version

To build the standalone version:

1. **Fetch standalone branch** - `rel-2.1-standalone`
   The standalone version contains a build profile to create a `.jar` file for analyzing single files.

2. **Build standalone analysis `.jar` file:**
   Run the following command:
   ```sh
   mvn clean package -Pstandalone -DskipTests
   ```
   
   The build will create a `.jar` file at:
   ```
   maven-plugin/target/fiori-js-analysis-standalone.jar
   ```

3. **Copy this file** to a permanent location, as new runs will clean the `/target` directory.

4. For more details, refer to `README-STANDALONE.md` in the standalone branch.

### Target Project Preparation

To run ESLint checks, you'll need to provide the path to the ESLint binary. We recommend adding ESLint as a dependency to your target project and using it in tests. Simply add a new dependency to `devDependencies`/`dependencies`:

```json
"eslint": "8.32.0"
```

Then run:
```sh
npm i
```

You can find the ESLint binary at:
```
node_modules/eslint/bin/eslint.js
```

For the testing utility, we recommend using the absolute path:
```
/Users/user/TestApps/ca.infra.testapp/node_modules/eslint/bin/eslint.js
```

## Test Execution

Navigate to `main_test.go` - it should already contain definitions for the executors:
- `analysisPluginA` - Pure analysis plugin version
- `analysisPluginB` - Version with disabled JSAnalysis checks
- `eslintPlugin` - ESLint custom plugin

Under the hood, executors simply run commands to execute code checks.

For the analysis plugin executor, you must provide the following values to configure it:
- `targetDir` - Directory of the project you want to check
- `version` - Version of the analysis plugin for this executor (ensure this version exists in your `.m2` repository)
- `standaloneFilePath` - Path to the standalone `.jar` file (leave as empty string if you don't want to run checks against individual files)

```go
// Pure Analysis Plugin
var analysisPluginA = executor.NewAnalysisPluginExecutor(
	targetDir,
	"2.1.8-A",
	"someDir/fiori-js-analysis-standalone-A.jar",
)
```

For the ESLint executor, you must provide the following:
- `targetDir` - Directory of the project you want to check
- `configFilePath` - `.eslintrc` configuration file path
- `rulesDir` - Path to custom ESLint rules
- `binaryPath` - Path to ESLint binary (`eslint.js` file)

```go
var eslintPlugin = executor.NewEsLintExecutor(
	targetDir,
	"/Users/user/eslint-plugin-fiori-custom/configure.eslintrc",
	"/Users/user/eslint-plugin-fiori-custom/lib/rules",
	"/Users/user/TestApps/ca.infra.testapp/node_modules/eslint/bin/eslint.js",
)
```

Then simply run the test suite:
```sh
go test -v -run TestOneRule -timeout 10m
```

**Note:** If you add more tests, you might need to increase the timeout.

## How to Add New Tests

To run the test suite, you need the following information:
- ESLint rule name
- Analysis plugin rule name
- Sample file with code that violates the expected check

The main test suite - `TestRules` - runs a single check against a sample file, then copies it to the target project and runs Analysis Plugin A/B checks and ESLint checks against the project.

The system then parses logs and returns any rule violations as errors. For tests to pass, we must ensure that our rule has been violated - the log entry with our rule must be present in the check report.

To add a test for a new rule, do the following:
- Add a new sample file to the `__test_files__` directory, then update:
1. `ruleNameToTestFile` - Use ESLint rule name as key and path to file as value
2. `eslintRuleNameToAnalysisName` - Use ESLint rule name as key and analysis plugin rule name as value

3. `RULES_TO_BE_TEST` - Add the ESLint rule name to this slice, so it will be included in the test suite

### Steps on a Linux machine

0. Store your token `git config credential.helper store`
1. clone the repo
2. Call `cd /.m2 (maven settings)
3. touch settings.xml
4. Set xml
`
<settings>
    <!--  Id: com.sap:artifactory:1.0.0:settings.xml  -->
    <mirrors>
        <mirror>
            <id>mirror1</id>
            <url>
                https://int.repositories.cloud.sap/artifactory/build-milestones/
            </url>
            <mirrorOf>*,!artifactory</mirrorOf>
        </mirror>
    </mirrors>
    <profiles>
        <profile>
            <id>release.build</id>
            <repositories>
                <repository>
                    <id>artifactory</id>
                    <url>
                        https://int.repositories.cloud.sap/artifactory/build-releases/
                    </url>
                </repository>
            </repositories>
            <properties>
                <tycho.disableP2Mirrors>true</tycho.disableP2Mirrors>
                <tycho.localArtifacts>ignore</tycho.localArtifacts>
            </properties>
        </profile>
        <profile>
            <id>milestone.build</id>
            <pluginRepositories>
                <pluginRepository>
                    <id>artifactory</id>
                    <url>
                        https://int.repositories.cloud.sap/artifactory/build-milestones/
                    </url>
                </pluginRepository>
            </pluginRepositories>
            <repositories>
                <repository>
                    <id>artifactory</id>
                    <url>
                        https://int.repositories.cloud.sap/artifactory/build-milestones/
                    </url>
                </repository>
            </repositories>
            <properties>
                <tycho.disableP2Mirrors>true</tycho.disableP2Mirrors>
                <tycho.localArtifacts>ignore</tycho.localArtifacts>
            </properties>
        </profile>
        <profile>
            <id>snapshot.build</id>
            <pluginRepositories>
                <pluginRepository>
                    <id>artifactory</id>
                    <url>
                        https://int.repositories.cloud.sap/artifactory/build-snapshots/
                    </url>
                </pluginRepository>
            </pluginRepositories>
            <repositories>
                <repository>
                    <id>artifactory</id>
                    <url>
                        https://int.repositories.cloud.sap/artifactory/build-snapshots/
                    </url>
                </repository>
            </repositories>
            <properties>
                <tycho.disableP2Mirrors>true</tycho.disableP2Mirrors>
                <tycho.localArtifacts>ignore</tycho.localArtifacts>
            </properties>
        </profile>
        <profile>
            <id>sonar</id>
            <activation>
                <activeByDefault>true</activeByDefault>
            </activation>
            <properties>
                <sonar.host.url>https://sonar.tools.sap</sonar.host.url>
                <tycho.disableP2Mirrors>true</tycho.disableP2Mirrors>
                <tycho.localArtifacts>ignore</tycho.localArtifacts>
            </properties>
        </profile>
    </profiles>
    <activeProfiles>
        <activeProfile>snapshot.build</activeProfile>
    </activeProfiles>
    <pluginGroups>
        <pluginGroup>com.sap.ldi</pluginGroup>
    </pluginGroups>
</settings>
`
5. For the repo run `mvn clean install`

6. switch to `IGNORErel2-1`

7. For this branch run `mvn clean install`

8.  get npm by sudo 





   







