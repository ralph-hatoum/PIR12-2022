/*
 * server.c
 * Version 20161003
 * Written by Harry Wong (RedAndBlueEraser)
 */

#include <netinet/in.h>
#include <pthread.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <unistd.h>
#include <string.h>

#define BACKLOG 10

typedef struct pthread_arg_t
{
    int new_socket_fd;
    struct sockaddr_in client_address;
    /* TODO: Put arguments passed to threads here. See lines 116 and 139. */
} pthread_arg_t;

/* Thread routine to serve connection to client. */
void *pthread_routine(void *arg);

/* Signal handler to handle SIGTERM and SIGINT signals. */
void signal_handler(int signal_number);

int socket_fd;

int main(int argc, char *argv[])
{
    int port, new_socket_fd;
    struct sockaddr_in address;
    pthread_attr_t pthread_attr;
    pthread_arg_t *pthread_arg;
    pthread_t pthread;
    socklen_t client_address_len;

    /* Get port from command line arguments or stdin. 
    port = argc > 1 ? atoi(argv[1]) : 0;
    if (!port)
    {
        printf("Enter Port: ");
        scanf("%d", &port);
    }
*/
    port = 8080;
    /* Initialise IPv4 address. */
    memset(&address, 0, sizeof address);
    address.sin_family = AF_INET;
    address.sin_port = htons(port);
    address.sin_addr.s_addr = INADDR_ANY;

    /* Create TCP socket. */
    if ((socket_fd = socket(AF_INET, SOCK_STREAM, 0)) == -1)
    {
        perror("socket");
        exit(1);
    }

    printf("%d\n", socket_fd);

    /* Bind address to socket. */
    if (bind(socket_fd, (struct sockaddr *)&address, sizeof address) == -1)
    {
        perror("bind");
        exit(1);
    }
    printf("Listening on port 8080\n");
    /* Listen on socket. */
    if (listen(socket_fd, BACKLOG) == -1)
    {
        perror("listen");
        exit(1);
    }

    /* Assign signal handlers to signals. */
    if (signal(SIGPIPE, SIG_IGN) == SIG_ERR)
    {
        perror("signal");
        exit(1);
    }
    if (signal(SIGTERM, signal_handler) == SIG_ERR)
    {
        perror("signal");
        exit(1);
    }
    if (signal(SIGINT, signal_handler) == SIG_ERR)
    {
        perror("signal");
        exit(1);
    }

    /* Initialise pthread attribute to create detached threads. */
    if (pthread_attr_init(&pthread_attr) != 0)
    {
        perror("pthread_attr_init");
        exit(1);
    }
    if (pthread_attr_setdetachstate(&pthread_attr, PTHREAD_CREATE_DETACHED) != 0)
    {
        perror("pthread_attr_setdetachstate");
        exit(1);
    }
    while (1)
    {

        /* Create pthread argument for each connection to client. */
        /* TODO: malloc'ing before accepting a connection causes only one small
         * memory when the program exits. It can be safely ignored.
         */
        pthread_arg = (pthread_arg_t *)malloc(sizeof *pthread_arg);
        if (!pthread_arg)
        {
            perror("malloc");
            continue;
        }

        /* Accept connection to client. */
        client_address_len = sizeof pthread_arg->client_address;
        new_socket_fd = accept(socket_fd, (struct sockaddr *)&pthread_arg->client_address, &client_address_len);
        if (new_socket_fd == -1)
        {
            perror("accept");
            free(pthread_arg);
            continue;
        }

        /* Initialise pthread argument. */
        pthread_arg->new_socket_fd = new_socket_fd;
        /* TODO: Initialise arguments passed to threads here. See lines 22 and
         * 139.
         */

        /* Create thread to serve connection to client. */
        if (pthread_create(&pthread, &pthread_attr, pthread_routine, (void *)pthread_arg) != 0)
        {
            perror("pthread_create");
            free(pthread_arg);
            continue;
        }
    }
    /* 
     * TODO: If you really want to close the socket, you would do it in
     * signal_handler(), meaning socket_fd would need to be a global variable.
     */
    signal_handler(SIGINT);

    return 0;
}

void *pthread_routine(void *arg)
{
    pthread_detach(pthread_self());
    pthread_arg_t *pthread_arg = (pthread_arg_t *)arg;
    int new_socket_fd = pthread_arg->new_socket_fd;
    struct sockaddr_in client_address = pthread_arg->client_address;
    /* TODO: Get arguments passed to threads here. See lines 22 and 116. */

    free(arg);

    /* TODO: Put client interaction code here. For example, use
     * write(new_socket_fd,,) and read(new_socket_fd,,) to send and receive
     * messages with the client.
     */
    char recv[8];
    read(new_socket_fd, &recv, 8);
    printf("%s\n", recv);

    char to_send[8];
    strcpy(to_send, "done");
    printf("%s\n", to_send);
    int resp;
    resp = send(new_socket_fd, &to_send, 8, 0);
    printf("%d\n", resp);
    close(new_socket_fd);
    return NULL;
}

void signal_handler(int signal_number)
{
    /* TODO: Put exit cleanup code here. */
    close(socket_fd);
    printf("Closed server");
    exit(0);
}