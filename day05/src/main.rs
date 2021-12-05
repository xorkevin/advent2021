use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

#[derive(Hash, PartialEq, Eq, Clone, Copy)]
struct Pos(i32, i32);

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut grid = HashMap::new();
    let mut grid2 = HashMap::new();
    for line in reader.lines() {
        if let [x1, y1, x2, y2] = line?
            .splitn(2, " -> ")
            .flat_map(|i| i.splitn(2, ","))
            .flat_map(|i| i.parse::<i32>())
            .collect::<Vec<_>>()[..]
        {
            if x1 == x2 {
                let (start, stop) = minmax(y1, y2);
                for j in start..=stop {
                    let k = Pos(x1, j);
                    *grid.entry(k).or_insert(0) += 1;
                    *grid2.entry(k).or_insert(0) += 1;
                }
            } else if y1 == y2 {
                let (start, stop) = minmax(x1, x2);
                for j in start..=stop {
                    let k = Pos(j, y1);
                    *grid.entry(k).or_insert(0) += 1;
                    *grid2.entry(k).or_insert(0) += 1;
                }
            } else {
                let mut start = Pos(x1, y1);
                let stop = Pos(x2, y2);
                let dir_x = sign(x2 - x1);
                let dir_y = sign(y2 - y1);
                while start.0 != stop.0 {
                    *grid2.entry(start).or_insert(0) += 1;
                    start = Pos(start.0 + dir_x, start.1 + dir_y);
                }
                *grid2.entry(stop).or_insert(0) += 1;
            }
        }
    }
    println!("Part 1: {}", grid.values().filter(|&&i| i > 1).count());
    println!("Part 2: {}", grid2.values().filter(|&&i| i > 1).count());
    Ok(())
}

fn minmax(a: i32, b: i32) -> (i32, i32) {
    if a < b {
        return (a, b);
    }
    return (b, a);
}

fn sign(a: i32) -> i32 {
    if a > 0 {
        1
    } else {
        -1
    }
}
