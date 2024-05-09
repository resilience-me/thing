#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>

void die(char *s) {
    perror(s);
    exit(1);
}

int main(void) {
    struct sockaddr_in si_other;
    int s, slen=sizeof(si_other);
    char buf[512];
    char message[512];

    if ((s=socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP)) == -1) {
        die("socket");
    }

    memset((char *) &si_other, 0, sizeof(si_other));
    si_other.sin_family = AF_INET;
    si_other.sin_port = htons(8888);

    if (inet_aton("172.234.123.103", &si_other.sin_addr) == 0) {
        fprintf(stderr, "inet_aton() failed\n");
        exit(1);
    }

    while(1) {
        printf("Enter message: ");
        if (fgets(message, sizeof(message), stdin) == NULL) break;  // Properly use fgets to avoid buffer overflow
        message[strcspn(message, "\n")] = 0;  // Remove newline character

        // Send the message
        if (sendto(s, message, strlen(message), 0 , (struct sockaddr *) &si_other, slen) == -1) {
            die("sendto()");
        }

        // Clear the buffer by filling null, it might have previously received data
        memset(buf, '\0', 512);
        // Try to receive some data, this is a blocking call
        if (recvfrom(s, buf, 512, 0, (struct sockaddr *) &si_other, &slen) == -1) {
            die("recvfrom()");
        }

        puts(buf);
    }

    close(s);
    return 0;
}
