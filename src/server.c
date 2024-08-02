#include <arpa/inet.h>
#include "handlers.h"

#define PORT 2012

CommandHandler command_handlers[256]  = {
    [0] = set_trustline
};

int main() {
    int sockfd;
    struct sockaddr_in server_addr, client_addr;
    socklen_t addr_len = sizeof(client_addr);
    Datagram dg;

    sockfd = socket(AF_INET, SOCK_DGRAM, 0);
    if (sockfd < 0) {
        return 1;
    }

    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(PORT);

    if (bind(sockfd, (const struct sockaddr *)&server_addr, sizeof(server_addr)) < 0) {
        return 1;
    }

    while (1) {
        int recv_len = recvfrom(sockfd, &dg, sizeof(dg), 0, (struct sockaddr *)&client_addr, &addr_len);
        if (recv_len < 0) {
            continue;
        }

        CommandHandler handler = command_handlers[(unsigned char)dg.command];
        if (handler) {
            handler(&dg, sockfd, &client_addr);
        }
    }
}
