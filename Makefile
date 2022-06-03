all: client_attack client_timeout server

client_attack: client_attack.c
	gcc client_attack.c -o client_attack

client_timeout: client_timeout.c 
	gcc client_timeout.c -o client_timeout

server: server.c 
	gcc server.c -o server

clean:
	rm -f client_attack client_timeout server