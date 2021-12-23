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

enum Pair {
    Val(i32),
    Pair((Box<Pair>, Box<Pair>)),
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

fn parse_tokens(tokens: &[Token]) -> Option<(Pair, &[Token])> {
    let (head, rest) = match tokens {
        [head, rest @ ..] => (head, rest),
        &[] => return None,
    };
    match head {
        Token::Num(val) => Some((Pair::Val(*val), rest)),
        Token::Lparen => {
            let (lhs, rest1) = match parse_tokens(rest) {
                Some((k, rest)) => (k, rest),
                None => return None,
            };
            let (rhs, rest2) = match parse_tokens(rest1) {
                Some((k, rest)) => (k, rest),
                None => return None,
            };
            Some((Pair::Pair((Box::new(lhs), Box::new(rhs))), rest2))
        }
        Token::Rparen => None,
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut root = None;
    for line in reader.lines() {
        let tokens = tokenize(line?.as_bytes().into_iter().peekable())?;
        let pair = match parse_tokens(&tokens) {
            Some((pair, _)) => pair,
            None => return Err("Invalid line".into()),
        };
        let k = match root {
            Some(r) => Pair::Pair((Box::new(r), Box::new(pair))),
            None => pair,
        };
        root = Some(k);
        println!("{:?}", tokens);
    }
    let mut _root = match root {
        Some(r) => r,
        None => return Err("Invalid line".into()),
    };
    Ok(())
}
