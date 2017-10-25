jaeger:
	docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p5775:5775/udp -p6831:6831/udp -p6832:6832/udp \
  -p5778:5778 -p16686:16686 -p14268:14268 -p9411:9411 jaegertracing/all-in-one:latest
consul:
	docker run -d -p 8400:8400 -p 8500:8500 -p 8600:53/udp -h node1 progrium/consul -server -bootstrap -ui-dir /ui

mysql:
	docker run --name gophercon-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=gopherconsql -d mysql:5

mysql-shell:
	docker run -it --link gophercon-mysql:mysql --rm mysql sh -c 'exec mysql -h"$$MYSQL_PORT_3306_TCP_ADDR" -P"$$MYSQL_PORT_3306_TCP_PORT" -uroot -p"$$MYSQL_ENV_MYSQL_ROOT_PASSWORD"'
mysql-create-config:
	docker run -it --link gophercon-mysql:mysql --rm mysql sh -c 'exec echo "create database config;" | exec mysql -h"$$MYSQL_PORT_3306_TCP_ADDR" -P"$$MYSQL_PORT_3306_TCP_PORT" -uroot -p"$$MYSQL_ENV_MYSQL_ROOT_PASSWORD"'

mysql-create-users:
	docker run -it --link gophercon-mysql:mysql --rm mysql sh -c 'exec echo "create database users;" | exec mysql -h"$$MYSQL_PORT_3306_TCP_ADDR" -P"$$MYSQL_PORT_3306_TCP_PORT" -uroot -p"$$MYSQL_ENV_MYSQL_ROOT_PASSWORD"'

mysql-create-profiles:
	docker run -it --link gophercon-mysql:mysql --rm mysql sh -c 'exec echo "create database profiles;" | exec mysql -h"$$MYSQL_PORT_3306_TCP_ADDR" -P"$$MYSQL_PORT_3306_TCP_PORT" -uroot -p"$$MYSQL_ENV_MYSQL_ROOT_PASSWORD"'

config-server:
	cd services/config-srv && go run main.go --database_url="root:gopherconsql@tcp(127.0.0.1:3306)/config"
