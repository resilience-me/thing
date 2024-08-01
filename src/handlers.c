void set_trustline(const Datagram *dg, int sockfd, struct sockaddr_in *client_addr) {
    char datadir[32];
    snprintf(datadir, sizeof(datadir), "%s/.ripple", getenv("HOME"));
    
    char peer[160];
    snprintf(peer, sizeof(peer), "%s/accounts/%s/peers/%s/%s", datadir, dg->x_username, dg->y_server_address, dg->y_username);

    if (access(peer, F_OK) == -1) {
        return; // Peer account does not exist
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
    memcpy(&counter, dg->counter, sizeof(counter));  // Copy the counter bytes
    counter = ntohl(counter);  // Convert from network to host byte order

    if (counter <= prev_counter) {
        return; // Counter already used or invalid
    }

    char trustline_out_path[192];
    snprintf(trustline_out_path, sizeof(trustline_out_path), "%s/trustline_out.txt", peer);

    FILE *trustline_file = fopen(trustline_out_path, "w");
    if (trustline_file) {
        int trustline_value = atoi(dg->arguments);
        fwrite(&trustline_value, sizeof(int), 1, trustline_file);
        fclose(trustline_file);
    }

    prev_counter = counter;  // Update prev_counter with the new value
    counter_file = fopen(counter_out_path, "w");
    if (counter_file) {
        prev_counter = htonl(prev_counter);  // Convert to network byte order before writing
        fwrite(&prev_counter, sizeof(int), 1, counter_file);
        fclose(counter_file);
    }
}
