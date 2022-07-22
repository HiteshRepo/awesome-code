### Steps To install mockery 
- this is global, not part of project dependency
- Repo: https://github.com/vektra/mockery
- brew install mockery
- brew upgrade mockery
  Or
- go install github.com/vektra/mockery/v2@latest

### Generate a mock
1. Create an interface for which you want a mock.
2. For example, if you want to test controller of a 3-tier app
   1. since controller is dependent of service and repo.
   2. mock is needed to be generated for service and repo in order to test controller logic only.
   3. so service and repo need to have interface of their own for mockery tool to generate mocks for them
3. command to generate mocks
   1. cd <directory where the interface is present>
   2. mockery --name=<interface_name>
   3. the mocks will be generated in a nested package.
   4. if the interface is accessible at github.com/repoName/projectName/packageName then mock is accessible at github.com/repoName/projectName/packageName/mocks