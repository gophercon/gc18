# Config Server

Config server is a microservice to store dynamic configuration. It can be used with the go-os.

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
	$ go get github.com/micro/config-srv
	$ config-srv --database_url="root:root@tcp(192.168.99.100:3306)/config"
	```

	OR as a docker container

	```shell
	$ docker run microhq/config-srv --database_url="root:root@tcp(192.168.99.100:3306)/config" --registry_address=YOUR_REGISTRY_ADDRESS
	```

## The API
Config server implements the following RPC Methods

Config
- Read
- Create
- Update
- Delete
- Search
- AuditLog
- Watch


### Config.Create
Create a record with id "NAMESPACE:CONFIG" at the path "supported_phones"
```shell
$ micro query go.micro.srv.config Config.Create '{
	"change": {
		"id": "NAMESPACE:CONFIG",
		"path": "supported_phones",
		"author": "asim",
		"comment": "adding ios phones", 
		"change_set": { "data": "{\"ios\": [\"4.0\", \"5s\", \"6\"]}"}
	}
}'

```

### Config.Read
Read record with id "NAMESPACE:CONFIG"
```shell
$ micro query go.micro.srv.config Config.Read '{"id": "NAMESPACE:CONFIG"}'

{
	"change": {
		"change_set": {
			"checksum": "edeea1cd6352f583",
			"data": "{\"supported_phones\":{\"ios\":[\"4.0\",\"5s\",\"6\"]}}",
			"source": "json",
			"timestamp": 1.452551133e+09
		},
		"id": "NAMESPACE:CONFIG"
	}
}
```

### Config.Update

Update record id "NAMESPACE:CONFIG" path "supported_phones/android" with android phones

```shell
$ micro query go.micro.srv.config Config.Update '{
	"change": {
		"id": "NAMESPACE:CONFIG",
		"path": "supported_phones/android",
		"author": "asim",
		"comment": "adding android phones",
		"change_set": {"data": "[\"galaxy nexus\", \"nexus 5\"]" }
	}
}'
```

### Config.Search

Search for configs
```shell
$ micro query go.micro.srv.config Config.Search

{
	"configs": [
		{
			"change_set": {
				"checksum": "37c9493aee25f277",
				"data": "{\"supported_phones\":{\"android\":[\"galaxy nexus\",\"nexus 5\"],\"ios\":[\"4.0\",\"5s\",\"6\"]}}",
				"source": "json",
				"timestamp": 1.452551694e+09
			},
			"id": "NAMESPACE:CONFIG"
		}
	]
}
```

### Config.Delete
Delete the ios config

```shell
$ micro query go.micro.srv.config Config.Delete '{
		"change": {
			"id": "NAMESPACE:CONFIG",
			"path": "supported_phones/ios",
			"author": "asim",
			"comment": "deleting ios config"
		}
}'
```

### Config.AuditLog

Read the audit log

```shell
$ micro query go.micro.srv.config Config.AuditLog

{
	"changes": [
		{
			"action": "create",
			"change": {
				"author": "asim",
				"change_set": {
					"checksum": "edeea1cd6352f583",
					"data": "{\"supported_phones\":{\"ios\":[\"4.0\",\"5s\",\"6\"]}}",
					"source": "json",
					"timestamp": 1.452551133e+09
				},
				"comment": "adding ios phones",
				"id": "NAMESPACE:CONFIG",
				"path": "supported_phones",
				"timestamp": 1.452551133e+09
			}
		},
		{
			"action": "update",
			"change": {
				"author": "asim",
				"change_set": {
					"checksum": "37c9493aee25f277",
					"data": "{\"supported_phones\":{\"android\":[\"galaxy nexus\",\"nexus 5\"],\"ios\":[\"4.0\",\"5s\",\"6\"]}}",
					"source": "json",
					"timestamp": 1.452551694e+09
				},
				"comment": "adding android phones",
				"id": "NAMESPACE:CONFIG",
				"path": "supported_phones/android",
				"timestamp": 1.452551694e+09
			}
		},
		{
			"action": "update",
			"change": {
				"author": "asim",
				"change_set": {
					"checksum": "47c3b86c7806648e",
					"data": "{\"supported_phones\":{\"android\":[\"galaxy nexus\",\"nexus 5\"]}}",
					"source": "json",
					"timestamp": 1.452551847e+09
				},
				"comment": "deleting ios config",
				"id": "NAMESPACE:CONFIG",
				"path": "supported_phones/ios",
				"timestamp": 1.452551847e+09
			}
		}
	]
}
```

### Config.Watch

Watch the config changes, meanwhile update the config in another terminal.

```shell
$ micro stream go.micro.srv.config Config.Watch '{"id": "NAMESPACE:CONFIG"}'

{
	"change_set": {
		"checksum": "95d998cdd93099c6",
		"data": "{\"supported_phones\":{\"android\":[\"galaxy nexus\",\"nexus 5\"]}}",
		"source": "json",
		"timestamp": 1.479574682e+09
	},
	"id": "NAMESPACE:CONFIG"
}

```

