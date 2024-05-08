#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>

#define BUFLEN 512  // Max length of buffer
#define PORT 8888   // The port on which to listen for incoming data

#define OPCODE_MESSAGE 1
#define OPCODE_GET_LAST_CMD 2
#define OPCODE_SET_TRUSTLINE 3
#define OPCODE_PAY 4

void die(char *s) {
    perror(s);
    exit(1);
}

void log_command(const char *command) {
    FILE *file = fopen("command_log.txt", "a");
    if (file == NULL) {
        die("Failed to open log file");
    }
    fprintf(file, "%s\n", command);
    fclose(file);
}

char* get_last_command() {
    static char last_command[BUFLEN];
    FILE *file = fopen("command_log.txt", "r");
    if (file == NULL) {
        return "No commands logged.";
    }
    last_command[0] = '\0'; // Clear the buffer
    char line[BUFLEN];
    while (fgets(line, BUFLEN, file) != NULL) {
        strcpy(last_command, line);
    }
    fclose(file);
    // Remove newline if present
    size_t len = strlen(last_command);
    if (len > 0 && last_command[len - 1] == '\n') {
        last_command[len - 1] = '\0';
    }
    return last_command;
}

typedef struct {
    int opcode;
    char data[BUFLEN];
} Command;

int main(void) {
    struct sockaddr_in si_me, si_other;
    
    int s, slen = sizeof(si_other), recv_len;
    char buf[BUFLEN];
    Command cmd;
    
    // Create a UDP socket
    if ((s = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP)) == -1) {
        die("socket");
    }
    
    // Zero out the structure
    memset((char *) &si_me, 0, sizeof(si_me));
    
    si_me.sin_family = AF_INET;
    si_me.sin_port = htons(PORT);
    si_me.sin_addr.s_addr = htonl(INADDR_ANY);
    
    // Bind socket to port
    if (bind(s, (struct sockaddr*)&si_me, sizeof(si_me)) == -1) {
        die("bind");
    }
    
    // Keep listening for data
    while (1) {
        printf("Waiting for data...\n");
        fflush(stdout);
        
        // Try to receive some data, this is a blocking call
        if ((recv_len = recvfrom(s, &cmd, sizeof(cmd), 0, (struct sockaddr *) &si_other, &slen)) == -1) {
            die("recvfrom()");
        }
        
        printf("Received packet from %s:%d\n", inet_ntoa(si_other.sin_addr), ntohs(si_other.sin_port));
        printf("Data: %s with opcode %d\n", cmd.data, cmd.opcode);
        
        switch (cmd.opcode) {
            case OPCODE_MESSAGE:
            case OPCODE_SET_TRUSTLINE:
            case OPCODE_PAY:
                log_command(cmd.data);
                sendto(s, "ACK", 3, 0, (struct sockaddr*)&si_other, slen);
                break;
            case OPCODE_GET_LAST_CMD:
                char *last_command = get_last_command();
                sendto(s, last_command, strlen(last_command), 0, (struct sockaddr*)&si_other, slen);
                break;
            default:
                sendto(s, "Invalid opcode", 14, 0, (struct sockaddr*)&si_other, slen);
                break;
        }
    }

    close(s);
    return 0;
}
