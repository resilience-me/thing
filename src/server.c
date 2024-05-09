#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>

#define PORT 2012
#define DATAGRAM_SIZE 389

int main() {
    int sockfd;
    struct sockaddr_in6 server_addr;
    struct sockaddr_storage client_addr;
    socklen_t addr_len = sizeof(client_addr);
    char buffer[DATAGRAM_SIZE];

    // Create socket
    if ((sockfd = socket(AF_INET6, SOCK_DGRAM, 0)) == -1) {
        perror("socket creation failed");
        exit(EXIT_FAILURE);
    }

    // Set socket option to allow both IPv4 and IPv6
    int opt = 0; // 0 enables dual-stack mode
    if (setsockopt(sockfd, IPPROTO_IPV6, IPV6_V6ONLY, (void *)&opt, sizeof(opt)) == -1) {
        perror("setsockopt failed");
        exit(EXIT_FAILURE);
    }

    // Initialize server address
    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin6_family = AF_INET6;
    server_addr.sin6_addr = in6addr_any;
    server_addr.sin6_port = htons(PORT);

    // Bind socket to both IPv4 and IPv6 addresses
    if (bind(sockfd, (struct sockaddr *)&server_addr, sizeof(server_addr)) == -1) {
        perror("bind failed");
        exit(EXIT_FAILURE);
    }

    // Server loop
    while (1) {
        // Receive data
        ssize_t recv_len = recvfrom(sockfd, buffer, BUFFER_SIZE, 0, (struct sockaddr*)&client_addr, &addr_len);
        if (recv_len == -1) {
            perror("recvfrom failed");
            continue;
        }

        // Deserialize received data
        Datagram dg;
        deserialize_datagram(buffer, &dg);

        // Verify signature and nonce
        if (!verify_signature(buffer, &dg) || !verify_nonce(dg.nonce)) continue;

        // Call appropriate command handler
        CommandHandler handler = command_handlers[dg.command];
        if (handler) {
            handler(&dg, sockfd, *(struct sockaddr_in *)&client_addr);
        }
    }

    close(sockfd);
    return 0;
}
