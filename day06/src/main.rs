use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut nums = vec![0; 9];
    for line in reader.lines() {
        for i in line?
            .split(",")
            .flat_map(|i| i.parse::<usize>())
            .filter(|&i| i < 9)
        {
            nums[i] += 1;
        }
    }
    for k in 0..256 {
        let mut next = vec![0; 9];
        for (n, i) in nums.iter().enumerate() {
            if n == 0 {
                next[6] += i;
                next[8] += i;
            } else {
                next[n - 1] += i;
            }
        }
        nums = next;
        if k == 79 {
            println!("Part 1: {}", nums.iter().sum::<i64>());
        }
    }
    println!("Part 2: {}", nums.iter().sum::<i64>());
    Ok(())
}
