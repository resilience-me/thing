func ToString(b []byte) string {
    return string(bytes.TrimRight(b, "\x00"))
}
