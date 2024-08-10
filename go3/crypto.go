// verifyHMAC checks the integrity and authenticity of the received buffer
func verifyHMAC(buf []byte, key []byte) bool {
    // The signature is the last 32 bytes of the buffer
    data := buf[:len(buf)-32]
    signature := buf[len(buf)-32:]

    mac := hmac.New(sha256.New, key)
    mac.Write(data)
    expectedMAC := mac.Sum(nil)

    return hmac.Equal(signature, expectedMAC)
}
