### Initialize a project
1. go mod init <module name>
2. ideal convention of naming a module is github_repo_url.
3. for example, in this case, it would be github.com/hiteshpattanayak-tw/golang-training

### Installing dependencies
1. go get <dependency_url>
2. for example: github.com/hiteshpattanayak-tw/golang-training
3. to target a specific branch: github.com/hiteshpattanayak-tw/golang-training@master
4. to target a specific version: github.com/hiteshpattanayak-tw/golang-training@v1.2.0
5. to target a specific commit: github.com/hiteshpattanayak-tw/golang-training@a1234d5
6. when a dependency is not directly used in any of the package, it is marked as "//indirect" in go.mod file
7. once you use the package and run 'go mod tidy', the indirect goes away

### Update dependencies
1. For minor version changes: go get -u <dependency_url>
2. For major version changes: 
   1. change the import url to v2 
   2. run 'go build'
   3. for example, change import "github.com/urfave/cli" to import "github.com/urfave/cli/v2" and run 'go build' in terminal

### Remove dependencies
1. Remove all references from the code.
2. Run 'go mod tidy'

### Vendoring dependencies
1. Act of making a copy of the third-party packages your project depends on and placing them in a vendor directory within your project.
2. Ensures the stability of your production builds without having to rely on external services.
3. But also increase package size.
4. If a package suddenly disappears from the internet, you are covered.
5. To vendor packages: 'go mod tidy' and 'go mod vendor'.
