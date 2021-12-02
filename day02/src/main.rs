use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut pos = 0;
    let mut depth = 0;
    let mut pos2 = 0;
    let mut depth2 = 0;
    let mut aim = 0;
    for line in reader.lines() {
        if let [dir, arg] = line?.splitn(2, ' ').collect::<Vec<_>>()[..] {
            let num = arg.parse::<i32>()?;
            match dir {
                "forward" => {
                    pos += num;
                    pos2 += num;
                    depth2 += aim * num;
                }
                "down" => {
                    depth += num;
                    aim += num;
                }
                "up" => {
                    depth -= num;
                    aim -= num;
                }
                _ => {}
            }
        }
    }
    println!("Part 1: {}", pos * depth);
    println!("Part 2: {}", pos2 * depth2);
    Ok(())
}
