/*
 * client.c
 * Version 20161003
 * Written by Harry Wong (RedAndBlueEraser)
 */

#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <unistd.h>

#define SERVER_NAME_LEN_MAX 255

int main(int argc, char *argv[])
{
    char server_name[SERVER_NAME_LEN_MAX + 1] = {0};
    int server_port, socket_fd;
    struct hostent *server_host;
    struct sockaddr_in server_address;
    /*
   
    if (argc > 1)
    {
       strncpy(server_name, argv[1], SERVER_NAME_LEN_MAX);
    }
    else
    {
        printf("Enter Server Name: ");
        scanf("%s", server_name);
    }
*/
    /* Get server port from command line arguments or stdin. */
    /*
    server_port = argc > 2 ? atoi(argv[2]) : 0;
    if (!server_port)
    {
        printf("Enter Port: ");
        scanf("%d", &server_port);
    }
    */
    strcpy(server_name, "10.0.1.2");
    server_port = 8080;
    /* Get server host from server name. */
    server_host = gethostbyname(server_name);

    /* Initialise IPv4 server address with server host. */
    memset(&server_address, 0, sizeof server_address);
    server_address.sin_family = AF_INET;
    server_address.sin_port = htons(server_port);
    memcpy(&server_address.sin_addr.s_addr, server_host->h_addr, server_host->h_length);
    int counter;
    int count = 0;
    while (1)
    {
        printf("Wave %d\n", count);
        printf("    \n");
        sleep(0.5);
        count += 1;
        for (counter = 0; counter < 25; counter++)
        {
            /* Create TCP socket. */
            if ((socket_fd = socket(AF_INET, SOCK_STREAM, 0)) == -1)
            {
                perror("socket");
                exit(1);
            }

            /* Connect to socket with server address. */
            if (connect(socket_fd, (struct sockaddr *)&server_address, sizeof server_address) == -1)
            {
                perror("connect");
                exit(1);
            }

            /* TODO: Put server interaction code here. For example, use
     * write(socket_fd,,) and read(socket_fd,,) to send and receive messages
     * with the client.
     */

            char msg[8] = "I am pi3";
            char *p;
            int resp;
            p = &msg;
            resp = send(socket_fd, p, 8, 0);
            printf("Sent %d bytes\n", resp);
            close(socket_fd);

            /*
        Partie du code pour recevoir une réponse ; si on le met, le client attend après chaque requête; mais du coup on bombarde bcp moins vite
        char rcv_buffer[8] = "";
        printf("%s\n", rcv_buffer);
        resp = recv(socket_fd, &rcv_buffer, 8, 0);
        printf("Received %d bytes, message (counter = %d) : \n", resp, counter);
        printf("%s\n", rcv_buffer);
        printf("  \n");*/
        }
        printf("   \n");
    }
    return 0;
}