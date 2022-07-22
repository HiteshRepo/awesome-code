##TAKEAWAYS FROM TEST CONTAINERS

[Github link]: https://github.com/testcontainers/testcontainers-go

### Helps us in creating container request [FROM Image directly]
```
request := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{pgPort + "/tcp"},
		WaitingFor:   wait.ForSQL(pgPort+"/tcp", "postgres", getDBUrl).Timeout(pgContainerSetupTimeout * time.Second),
		BindMounts:   map[string]string{"/docker-entrypoint-initdb.d": pathToSeedSql},
		Env: map[string]string{
			"POSTGRES_DB":       pgDatabase,
			"POSTGRES_USER":     pgUsername,
			"POSTGRES_PASSWORD": pgPassword,
		},
		Networks:     []string{p.network.Name},
	}
	container, err := testcontainers.GenericContainer(p.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
	})
```

### Helps us in creating container request [By Building from an existing Dockerfile]
```
request := testcontainers.ContainerRequest{
		Name: "new_books_api",
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       suite.Config.RootDir,
			Dockerfile:    "build/package/Dockerfile",
			PrintBuildLog: true,
		},
		Env: map[string]string{
			"APP_PORT":    suite.appPort,
			"DB_USER":     pgUsername,
			"DB_NAME":     pgDatabase,
			"DB_PASSWORD": pgPassword,
			"DB_HOST":     suite.dbHost,
			"DB_PORT":     suite.dbPort,
		},
		ExposedPorts: []string{"5500"},
		Networks:     []string{suite.Network.Name},
		Cmd:          []string{"/bin/sh", "-c", "/wait && ./new_books_api -configFile default.yaml"},
		WaitingFor:   wait.ForLog("Starting server"),
	}

	app, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
	})
```

### Helps us start the container [using the container request]
```
err = suite.appContainer.Start(suite.ctx)
suite.Require().NoError(err)
```

### Helps us in creating the network
```
network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Name: clusterNetworkName},
})
```

### Helps us in removing the container after integration tests
```
if container != nil {
    if err := container.Terminate(p.ctx); err != nil {
        p.Require().NoError(err)
    }
}
```

### Helps us in removing the network after integration tests and removing containers
```
if f.Network != nil {
		_ = f.Network.Remove(f.Ctx)
}
```