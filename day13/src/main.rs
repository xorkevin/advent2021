use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

#[derive(PartialEq, Eq, Hash, Clone, Copy)]
struct Pos(usize, usize);

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut points = HashSet::new();
    let mut track_points = true;
    let mut first = true;
    for line in reader.lines() {
        let line = line?;
        if line == "" {
            track_points = false;
            continue;
        }
        if track_points {
            if let Some((xs, ys)) = line.split_once(",") {
                points.insert(Pos(xs.parse::<usize>()?, ys.parse::<usize>()?));
            } else {
                return Err("Invalid line".into());
            }
            continue;
        }
        let eqn = match line.rsplit_once(" ") {
            Some((_, eqn)) => eqn,
            _ => return Err("Invalid line".into()),
        };
        let (yaxis, val) = match eqn.split_once("=") {
            Some((axis, val)) => (axis == "y", val.parse::<usize>()?),
            _ => return Err("Invalid line".into()),
        };
        points = fold(points, yaxis, val);
        if first {
            first = false;
            println!("Part 1: {}", points.len());
        }
    }
    let (mx, my) = find_max(&points);
    let mut grid = (0..=my)
        .into_iter()
        .map(|_| vec![b' '; mx + 1])
        .collect::<Vec<_>>();
    for Pos(x, y) in points {
        grid[y][x] = b'#';
    }

    println!("Part 2:");
    for i in grid {
        println!("{}", String::from_utf8(i)?);
    }
    Ok(())
}

fn fold(points: HashSet<Pos>, yaxis: bool, axis: usize) -> HashSet<Pos> {
    let mut next = HashSet::new();
    for Pos(x, y) in points {
        if yaxis {
            if y <= axis {
                next.insert(Pos(x, y));
                continue;
            }
            next.insert(Pos(x, axis - (y - axis)));
        } else {
            if x <= axis {
                next.insert(Pos(x, y));
                continue;
            }
            next.insert(Pos(axis - (x - axis), y));
        }
    }
    next
}

fn find_max(points: &HashSet<Pos>) -> (usize, usize) {
    points.iter().fold((0, 0), |(mx, my), &Pos(x, y)| {
        (if x > mx { x } else { mx }, if y > my { y } else { my })
    })
}
