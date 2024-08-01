#include <arpa/inet.h>

#define PORT 2012

typedef struct {
    char command;
    char x_username[32];
    char y_username[32];
    char y_server_address[32];
    char arguments[256];
    char counter[4];
    char signature[32];
} Datagram;

typedef void (*CommandHandler)(Datagram*, int sockfd);

int main() {
    int sockfd;
    struct sockaddr_in server_addr, client_addr;
    socklen_t addr_len = sizeof(client_addr);
    Datagram datagram;

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
        int recv_len = recvfrom(sockfd, &datagram, sizeof(datagram), 0, (struct sockaddr *)&client_addr, &addr_len);
        if (recv_len < 0) {
            continue;
        }
    }
}
