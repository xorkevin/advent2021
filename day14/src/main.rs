use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut first = 0;
    let mut last = 0;
    let mut pairs = HashMap::new();
    let mut rules = HashMap::new();
    let mut track_pairs = true;
    for line in reader.lines() {
        let line = line?;
        if line == "" {
            track_pairs = false;
            continue;
        }
        if track_pairs {
            let v = line.as_bytes().into_iter().collect::<Vec<_>>();
            for k in v.windows(2) {
                let (a, b) = match k {
                    &[&a, &b] => (a, b),
                    _ => return Err("Invalid line".into()),
                };
                *pairs.entry((a, b)).or_insert(0) += 1;
            }
            first = *v[0];
            last = *v[v.len() - 1];
            continue;
        }
        let (lhs, rhs) = match line.split_once(" -> ") {
            Some(eqn) => eqn,
            _ => return Err("Invalid line".into()),
        };
        if let [&a, &b] = lhs.as_bytes().into_iter().collect::<Vec<_>>()[..] {
            rules.insert((a, b), rhs.as_bytes()[0]);
        } else {
            return Err("Invalid line".into());
        }
    }

    for _ in 0..10 {
        pairs = process_step(&rules, pairs);
    }
    let (max, min) = find_max_min(&pairs, first, last);
    println!("Part 1: {}", max - min);
    for _ in 0..30 {
        pairs = process_step(&rules, pairs);
    }
    let (max, min) = find_max_min(&pairs, first, last);
    println!("Part 2: {}", max - min);
    Ok(())
}

fn process_step(
    rules: &HashMap<(u8, u8), u8>,
    pairs: HashMap<(u8, u8), usize>,
) -> HashMap<(u8, u8), usize> {
    let mut next = HashMap::new();
    for ((a, b), v) in pairs {
        let c = rules[&(a, b)];
        *next.entry((a, c)).or_insert(0) += v;
        *next.entry((c, b)).or_insert(0) += v;
    }
    next
}

fn find_max_min(pairs: &HashMap<(u8, u8), usize>, first: u8, last: u8) -> (usize, usize) {
    let mut counts = HashMap::from([(first, 1), (last, 1)]);
    let mut max = 0;
    let mut min = 0;
    for (&(a, b), v) in pairs {
        *counts.entry(a).or_insert(0) += v;
        *counts.entry(b).or_insert(0) += v;
        max = counts[&b] / 2;
        min = counts[&b] / 2;
    }
    counts.into_iter().fold((max, min), |(max, min), (_, v)| {
        (
            if v / 2 > max { v / 2 } else { max },
            if v / 2 < min { v / 2 } else { min },
        )
    })
}
