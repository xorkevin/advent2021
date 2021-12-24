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

#[derive(Clone)]
enum Pair {
    Val(i32),
    Pair((Box<Pair>, Box<Pair>)),
}

impl Pair {
    fn add_left(&mut self, v: i32) {
        match self {
            Pair::Pair((lhs, _)) => lhs.add_left(v),
            Pair::Val(val) => *val += v,
        }
    }

    fn add_right(&mut self, v: i32) {
        match self {
            Pair::Pair((_, rhs)) => rhs.add_right(v),
            Pair::Val(val) => *val += v,
        }
    }

    fn explode(&mut self, depth: usize) -> Option<(i32, i32)> {
        let (lhs, rhs) = match self {
            Pair::Pair((lhs, rhs)) => (lhs, rhs),
            Pair::Val(_) => return None,
        };
        if depth > 3 {
            let lhs = match **lhs {
                Pair::Pair(_) => return None,
                Pair::Val(v) => v,
            };
            let rhs = match **rhs {
                Pair::Pair(_) => return None,
                Pair::Val(v) => v,
            };
            *self = Pair::Val(0);
            Some((lhs, rhs))
        } else {
            if let Some((l, r)) = lhs.explode(depth + 1) {
                if r != 0 {
                    rhs.add_left(r);
                }
                return Some((l, 0));
            }
            if let Some((l, r)) = rhs.explode(depth + 1) {
                if l != 0 {
                    lhs.add_right(l);
                }
                return Some((0, r));
            }
            None
        }
    }

    fn split(&mut self) -> bool {
        match self {
            Pair::Pair((lhs, rhs)) => lhs.split() || rhs.split(),
            Pair::Val(v) => {
                return if *v > 9 {
                    let l = *v / 2;
                    let r = *v - l;
                    *self = Pair::Pair((Box::new(Pair::Val(l)), Box::new(Pair::Val(r))));
                    true
                } else {
                    false
                }
            }
        }
    }

    fn reduce_step(&mut self) -> bool {
        if let Some(_) = self.explode(0) {
            return true;
        }
        return self.split();
    }

    fn reduce(&mut self) {
        while self.reduce_step() {}
    }

    fn magnitude(&self) -> i32 {
        match self {
            Pair::Pair((lhs, rhs)) => lhs.magnitude() * 3 + rhs.magnitude() * 2,
            Pair::Val(v) => *v,
        }
    }
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
            let rest3 = match rest2 {
                [Token::Rparen, rest @ ..] => rest,
                _ => return None,
            };
            Some((Pair::Pair((Box::new(lhs), Box::new(rhs))), rest3))
        }
        Token::Rparen => None,
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut root = None;
    let mut nums = Vec::new();
    for line in reader.lines() {
        let tokens = tokenize(line?.as_bytes().into_iter().peekable())?;
        let pair = match parse_tokens(&tokens) {
            Some((pair, _)) => pair,
            None => return Err("Invalid line".into()),
        };
        nums.push(pair.clone());
        let mut k = match root {
            Some(r) => Pair::Pair((Box::new(r), Box::new(pair))),
            None => pair,
        };
        k.reduce();
        root = Some(k);
    }
    match root {
        Some(r) => println!("Part 1: {}", r.magnitude()),
        None => return Err("No vals".into()),
    };

    match (0..nums.len() - 1)
        .into_iter()
        .flat_map(|i| {
            (i + 1..nums.len())
                .into_iter()
                .flat_map(move |j| vec![(i, j), (j, i)])
        })
        .map(|(i, j)| {
            let mut k = Pair::Pair((Box::new(nums[i].clone()), Box::new(nums[j].clone())));
            k.reduce();
            k.magnitude()
        })
        .max()
    {
        Some(v) => println!("Part 2: {}", v),
        None => return Err("No vals".into()),
    }

    Ok(())
}
