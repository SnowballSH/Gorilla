function Grammar() {
    return {
        Magic: 0x69,

        Integer: 0x00
        // Length of integer in unsigned leb128
        // Followed by leb128
    };
}

console.log(Grammar());