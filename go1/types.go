
// ToString converts a byte slice to a string, trimming trailing zeroes.
func ToString(b []byte) string {
    return string(bytes.TrimRight(b, "\x00"))
}
