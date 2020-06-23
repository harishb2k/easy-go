package easy

type Error struct {
    error
    Name        string
    Description string
    Err         error
}
