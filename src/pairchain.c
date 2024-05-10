void saveTransactionChain(const char *filename, TransactionNode *head) {
    FILE *file = fopen(filename, "wb");
    if (!file) {
        perror("Failed to open file for writing");
        return;
    }

    TransactionNode *current = head;
    while (current != NULL) {
        fwrite(&current->transaction, sizeof(Transaction), 1, file);
        current = current->next;
    }

    fclose(file);
}

void appendTransaction(const char *filename, Transaction *transaction) {
    FILE *file = fopen(filename, "ab");  // Open in append mode
    if (!file) {
        perror("Failed to open file for appending");
        return;
    }

    fwrite(transaction, sizeof(Transaction), 1, file);
    fclose(file);
}

TransactionNode *loadTransactionChain(const char *filename) {
    FILE *file = fopen(filename, "rb");
    if (!file) {
        perror("Failed to open file for reading");
        return NULL;
    }

    TransactionNode *head = NULL;
    TransactionNode *current = NULL;

    Transaction temp;
    while (fread(&temp, sizeof(Transaction), 1, file)) {
        TransactionNode *newNode = (TransactionNode *)malloc(sizeof(TransactionNode));
        if (!newNode) {
            perror("Failed to allocate memory for transaction node");
            fclose(file);
            return head;
        }
        newNode->transaction = temp;
        newNode->next = NULL;

        if (head == NULL) {
            head = newNode;
            current = head;
        } else {
            current->next = newNode;
            current = newNode;
        }
    }

    fclose(file);
    return head;
}
