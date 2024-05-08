#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <sys/socket.h>

#define SERVER "172.234.123.103"
#define BUFLEN 512  // Max length of buffer
#define PORT 8888   // The port on which to send data

#define OPCODE_MESSAGE 1
#define OPCODE_GET_LAST_CMD 2
#define OPCODE_SET_TRUSTLINE 3
#define OPCODE_PAY 4

void die(char *s) {
    perror(s);
    exit(1);
}

typedef struct {
    int opcode;
    char data[BUFLEN];
} Command;

int main(void) {
    struct sockaddr_in si_other;
    int s, slen = sizeof(si_other);
    Command cmd;

    if ((s = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP)) == -1) {
        die("socket");
    }

    memset((char *) &si_other, 0, sizeof(si_other));
    si_other.sin_family = AF_INET;
    si_other.sin_port = htons(PORT);

    if (inet_aton(SERVER, &si_other.sin_addr) == 0) {
        fprintf(stderr, "inet_aton() failed\n");
        exit(1);
    }

    while (1) {
        printf("Enter opcode (1-Message, 2-Get Last Cmd, 3-Set Trustline, 4-Pay): ");
        scanf("%d", &cmd.opcode);
        getchar(); // clear the newline character after the number

        if (cmd.opcode != OPCODE_GET_LAST_CMD) {
            printf("Enter message: ");
            fgets(cmd.data, BUFLEN, stdin);
            // Remove newline at the end of input if present
            size_t len = strlen(cmd.data);
            if (len > 0 && cmd.data[len-1] == '\n') {
                cmd.data[len-1] = '\0';
            }
        } else {
            strcpy(cmd.data, ""); // Send empty data for GET_LAST_CMD
        }

        // Send the command
        if (sendto(s, &cmd, sizeof(cmd), 0, (struct sockaddr *) &si_other, slen) == -1) {
            die("sendto()");
        }

        // Receive the response
        if (recvfrom(s, cmd.data, BUFLEN, 0, (struct sockaddr *) &si_other, &slen) == -1) {
            die("recvfrom()");
        }

        printf("Server reply: %s\n", cmd.data);
    }

    close(s);
    return 0;
}
