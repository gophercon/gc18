# Profile Server

Profile server is a service to store user profiles.

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Start a mysql database

4. Download and start the service

	```shell
	go get github.com/micro/profile-srv
	profile-srv --database_url="root:root@tcp(192.168.99.100:3306)/profile"
	```

	OR as a docker container

	```shell
	docker run microhq/profile-srv --database_url="root:root@tcp(192.168.99.100:3306)/profile" --registry_address=YOUR_REGISTRY_ADDRESS
	```

## The API
Profile server implements the following RPC Methods

Record
- Create
- Read
- Update
- Delete
- Search

### Record.Create
```shell
micro query go.micro.srv.profile Record.Create '{
	"profile": {
		"id": "1", 
		"name": "john", 
		"owner": "john", 
		"type": 0, 
		"display_name": "John Smith", 
		"blurb": "hey Im john, this is my profile", 
		"url": "http://example.com", 
		"location": "london"
	}
}' 
{}
```

### Record.Search
```shell
micro query go.micro.srv.profile Record.Search '{"limit": 10}'
{
	"profiles": [
		{
			"blurb": "hey Im john, this is my profile",
			"created": 1.453408904e+09,
			"id": "1",
			"location": "london",
			"name": "john",
			"owner": "john",
			"updated": 1.453408904e+09,
			"url": "http://example.com"
		}
	]
}
```
