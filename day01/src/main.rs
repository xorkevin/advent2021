use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let nums = {
        let mut nums = Vec::new();
        for line in reader.lines() {
            nums.push(line?.parse::<i32>()?);
        }
        nums
    };

    let mut count = 0;
    for [prev, next] in nums.windows(2).flat_map(<&[i32; 2]>::try_from) {
        if next > prev {
            count += 1;
        }
    }
    let mut count2 = 0;
    for [prev, next] in nums
        .windows(3)
        .flat_map(<&[i32; 3]>::try_from)
        .map(|[one, two, three]| one + two + three)
        .collect::<Vec<i32>>()
        .windows(2)
        .flat_map(<&[i32; 2]>::try_from)
    {
        if next > prev {
            count2 += 1;
        }
    }
    println!("Part 1: {}", count);
    println!("Part 2: {}", count2);
    Ok(())
}
