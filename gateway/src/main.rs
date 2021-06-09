fn main() {
    let mut m: u8 = 127;

    m = u8::wrapping_add(m, 130);

    println!("Wrapped {}", m);
}
