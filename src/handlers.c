#include <stdio.h>
#include <string.h>
#include <arpa/inet.h>
#include "mac/sha2.h"

void set_trustline(const Datagram *dg, int sockfd, struct sockaddr_in *client_addr) {
    char datadir[32];
    snprintf(datadir, sizeof(datadir), "%s/.ripple", getenv("HOME"));

    char peer[160];
    snprintf(peer, sizeof(peer), "%s/accounts/%s/peers/%s/%s", datadir, dg->x_username, dg->y_server_address, dg->y_username);
    
    if (access(peer, F_OK) == -1) {
        return;
    }

    char secret_key_path[192];
    snprintf(secret_key_path, sizeof(secret_key_path), "%s/secretkey.txt", peer);

    char secret_key[32];
    FILE *key_file = fopen(secret_key_path, "r");
    if (key_file) {
        fread(secret_key, sizeof(char), 32, key_file);
        fclose(key_file);
    }

    char data_with_key[sizeof(Datagram)];
    memcpy(data_with_key, dg, sizeof(Datagram) - sizeof(dg->signature));
    memcpy(data_with_key + sizeof(Datagram) - sizeof(dg->signature), secret_key, 32);

    char hash[32];  
    sha256(data_with_key, sizeof(data_with_key), hash);

    if (memcmp(hash, dg->signature, sizeof(dg->signature)) != 0) {
        return;
    }

    char counter_out_path[192];
    snprintf(counter_out_path, sizeof(counter_out_path), "%s/counter_out.txt", peer);

    int prev_counter;
    FILE *counter_file = fopen(counter_out_path, "r");
    if (counter_file) {
        fread(&prev_counter, sizeof(int), 1, counter_file);
        fclose(counter_file);
    }

    int counter;
    memcpy(&counter, dg->counter, sizeof(counter));
    counter = ntohl(counter);

    if (counter <= prev_counter) {
        return;
    }

    int trustline;
    memcpy(&trustline, dg->arguments, sizeof(trustline));
    trustline = ntohl(trustline);

    char trustline_out_path[192];
    snprintf(trustline_out_path, sizeof(trustline_out_path), "%s/trustline_out.txt", peer);
    
    FILE *trustline_file = fopen(trustline_out_path, "w");
    if (trustline_file) {
        fwrite(&trustline, sizeof(int), 1, trustline_file);  
        fclose(trustline_file);
    }

    counter++;
    FILE *counter_file_write = fopen(counter_out_path, "w");
    if (counter_file_write) {
        int network_counter = htonl(counter);
        fwrite(&network_counter, sizeof(int), 1, counter_file_write);
        fclose(counter_file_write);
    }
}
