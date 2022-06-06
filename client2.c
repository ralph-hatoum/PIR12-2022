#include <sys/socket.h>
#include <fcntl.h>
#include <errno.h>
#include <poll.h>
#include <time.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <unistd.h>
#include <arpa/inet.h>

int main(int argc, char **argv) {
    u_short port;                /* user specified port number */
    char *addr;                  /* will be a pointer to the address */
    struct sockaddr_in address;  /* the libc network address data structure */
    short int sock = -1;         /* file descriptor for the network socket */
    fd_set fdset;
    struct timeval tv;

    if (argc != 3) {
        fprintf(stderr, "Usage %s <port_num> <address>\n", argv[0]);
        return EXIT_FAILURE;
    }

    port = atoi(argv[1]);
    addr = argv[2];

    address.sin_family = AF_INET;
    address.sin_addr.s_addr = inet_addr(addr); /* assign the address */
    address.sin_port = htons(port);            /* translate int2port num */

    FD_ZERO(&fdset);
    FD_SET(sock, &fdset);
    tv.tv_sec = 5;             /* 10 second timeout */
    tv.tv_usec = 0;

    int cpt = 0;
    int v;
    for(v=0; v<100;v++){
        sock = socket(AF_INET, SOCK_STREAM, 0);
        fcntl(sock, F_SETFL, O_NONBLOCK);

        connect(sock, (struct sockaddr *)&address, sizeof(address));

        if (select(sock + 1, NULL, &fdset, NULL, &tv) == 1)
        {
            int so_error;
            socklen_t len = sizeof so_error;

            getsockopt(sock, SOL_SOCKET, SO_ERROR, &so_error, &len);

            if (so_error == 0) {
                printf("%s:%d is open\n", addr, port);
            }else{
                cpt++;
            }
        }

        close(sock);
        sock = -1;
    }
    printf("nombre de non recu sur 100 : %d\n",v);
    return 0;
}
