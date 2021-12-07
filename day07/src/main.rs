use std::cmp::min;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut nums = Vec::new();
    for line in reader.lines() {
        nums = line?
            .split(",")
            .flat_map(|i| i.parse::<i32>())
            .collect::<Vec<_>>();
    }
    nums.sort();
    let median = nums[nums.len() / 2];
    let (a1, a2) = {
        let k = nums.iter().sum::<i32>();
        let a = (k as f64) / nums.len() as f64;
        (a.floor() as i32, a.ceil() as i32)
    };
    let (d, d1, d2) = nums.iter().fold((0, 0, 0), |(d, d1, d2), i| {
        let k1 = (i - a1).abs();
        let k2 = (i - a2).abs();
        (
            d + (i - median).abs(),
            d1 + k1 * (k1 + 1) / 2,
            d2 + k2 * (k2 + 1) / 2,
        )
    });
    println!("Part 1: {}", d);
    println!("Part 2: {}", min(d1, d2));
    Ok(())
}
