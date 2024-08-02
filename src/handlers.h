#ifndef HANDLERS_H
#define HANDLERS_H

typedef struct {
    char command;
    char x_username[32];
    char y_username[32];
    char y_server_address[32];
    char arguments[256];
    char counter[4];
    char signature[32];
} Datagram;

typedef void (*CommandHandler)(const Datagram*, int, struct sockaddr_in*);

void set_trustline(const Datagram *dg, int sockfd, struct sockaddr_in *client_addr);

#endif /* HANDLERS_H */
