one:
	@echo ''
	@echo ''
	@cat client/functions/func.js
	@echo ''
	@echo ''
	@read -n 1 -s -r -p "Press any key to continue"
	@./client/build/fx up client/functions/func.js

multiple:
	@echo '#Node'
	@echo '=================================='
	@echo ''
	@cat client/functions/func.js
	@echo ''
	@echo ''
	@echo '#Golang'
	@echo '=================================='
	@echo ''
	@cat client/functions/func.go
	@echo ''
	@echo ''
	@echo '#Ruby'
	@echo '=================================='
	@echo ''
	@cat client/functions/func.rb
	@echo ''
	@echo ''
	@echo '#Python'
	@echo '=================================='
	@echo ''
	@cat client/functions/func.py
	@echo ''
	@echo ''
	@read -n 1 -s -r -p "Press any key to continue"
	@./client/build/fx up client/functions/func.js
alot:
	./100.sh


