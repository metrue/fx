up:
	fx up -n ping.go.0 -p 12341 --force fx.go
	fx up -n ping.java.0 -p 12342 --force fx.java
	fx up -n ping.js.0 -p 12343 --force fx.js 
	fx up -n ping.pl.0 -p 12344 --force fx.pl
	fx up -n ping.py.0 -p 12345 --force fx.py
	fx up -n ping.rb.0 -p 12346 --force fx.rb
	fx up -n ping.rs.0 -p 12347 --force fx.rs
	fx up -n ping.cr.0 -p 12348 --force fx.cr
	fx up -n ping.d.0 -p 12349 --force fx.d
down:
	fx down ping.go.0
	fx down ping.java.0
	fx down ping.js.0 
	fx down ping.pl.0 
	fx down ping.py.0 
	fx down ping.rb.0 
	fx down ping.rs.0 
	fx down ping.cr.0 
	fx down ping.d.0 
l:
	fx list
t:
	@curl 127.0.0.1:12341 
	@echo ""
	@curl 127.0.0.1:12343
	@echo ""
	@curl 127.0.0.1:12344
	@echo ""
	@curl 127.0.0.1:12345
	@echo ""
	@curl 127.0.0.1:12346
	@echo ""
	@curl 127.0.0.1:12348
	@echo ""
	@curl 127.0.0.1:12349
	@echo ""
	@curl -X POST -H 'Content-Type: application/json' 127.0.0.1:12342 -d '{"a": 1, "b": 1 }'
	@echo ""
	@curl -X POST -H 'Content-Type: application/json' 127.0.0.1:12347 -d '{"a": 1, "b": 2 }'
	@echo ""
	@curl -X POST -H 'Content-Type: application/json' 127.0.0.1:12349 -d '{"a": 1, "b": 3 }'
	@echo ""
b:
	ab -n 10000 -c 10	http://127.0.0.1:12343/
m:
	open http://127.0.0.1:12343/swagger-stats/ux
	open http://127.0.0.1:12343/swagger-stats/metrics
