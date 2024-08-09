
// ToString converts a [32]byte array to a string, trimming trailing zeroes.
func ToString(b []byte) string {
    return string(bytes.TrimRight(b, "\x00"))
}
