use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let scores = reader
        .lines()
        .flat_map(|i| i)
        .map(|i| parse_pair(i.as_bytes()))
        .collect::<Vec<_>>();
    println!("Part 1: {}", scores.iter().map(|&(_, i, _)| i).sum::<i32>());
    let mut completes = scores
        .into_iter()
        .filter_map(|(r, _, i)| if r == 1 { Some(i) } else { None })
        .collect::<Vec<_>>();
    completes.sort();
    if completes.len() > 0 {
        println!("Part 2: {}", completes[completes.len() / 2]);
    }
    Ok(())
}

fn parse_pair(s: &[u8]) -> (i32, i32, i64) {
    let mut stack = Vec::with_capacity(s.len());
    for &i in s {
        match i {
            b'(' | b'[' | b'{' | b'<' => stack.push(i),
            b')' | b']' | b'}' | b'>' => {
                if let Some(c) = stack.pop() {
                    match i {
                        b')' => {
                            if c != b'(' {
                                return (2, 3, 0);
                            }
                        }
                        b']' => {
                            if c != b'[' {
                                return (2, 57, 0);
                            }
                        }
                        b'}' => {
                            if c != b'{' {
                                return (2, 1197, 0);
                            }
                        }
                        b'>' => {
                            if c != b'<' {
                                return (2, 25137, 0);
                            }
                        }
                        _ => (),
                    }
                } else {
                    return (3, 0, 0);
                }
            }
            _ => return (4, 0, 0),
        }
    }
    if stack.len() != 0 {
        return (1, 0, translate(&stack));
    }
    (0, 0, 0)
}

fn translate(s: &[u8]) -> i64 {
    s.iter().rfold(0, |a, i| match i {
        b'(' => a * 5 + 1,
        b'[' => a * 5 + 2,
        b'{' => a * 5 + 3,
        b'<' => a * 5 + 4,
        _ => -1,
    })
}
