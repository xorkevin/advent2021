use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::iter::Peekable;
use std::slice::Iter;

const PUZZLEINPUT: &str = "input.txt";

#[derive(Debug)]
enum Token {
    Lparen,
    Rparen,
    Num(i32),
}

fn tokenize(mut b: Peekable<Iter<u8>>) -> Result<Vec<Token>, Box<dyn std::error::Error>> {
    let mut tokens = Vec::new();
    while let Some(&c) = b.next() {
        match c {
            b'[' => tokens.push(Token::Lparen),
            b']' => tokens.push(Token::Rparen),
            b'0'..=b'9' => {
                let mut buf = vec![c];
                while let Some(&c2) = b.next_if(|&&i| i >= b'0' && i <= b'9') {
                    buf.push(c2);
                }
                tokens.push(Token::Num(String::from_utf8(buf)?.parse::<i32>()?));
            }
            _ => (),
        }
    }
    Ok(tokens)
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    for line in reader.lines() {
        let tokens = tokenize(line?.as_bytes().into_iter().peekable())?;
        println!("{:?}", tokens);
    }
    Ok(())
}
